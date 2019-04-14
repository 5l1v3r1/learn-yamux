package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	pb "github.com/jimmy-xu/learn-yamux/protocols/grpc"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

var defaultDialTimeout = 15 * time.Second
var defaultCloseTimeout = 5 * time.Second

const (
	unixSocketScheme  = "unix"
)

// AgentClient is an agent gRPC client connection wrapper for agentgrpc.AgentServiceClient
type AgentClient struct {
	pb.GreeterClient
	conn *grpc.ClientConn
}


func main() {
	ctx := context.Background()
	socks := ""
	if len(os.Args) > 1 {
		socks = os.Args[1]
	}
	if socks == "" {
		logrus.Fatalf("please specify socks file, for example ~/Documents/antfin/serial/com1")
	}

	client, err := NewAgentClient(ctx, socks, true)
	if err != nil {
		logrus.Fatalf("NewAgentClient error:%v", err)
	}

	in := pb.HelloRequest{
		Name: "world",
	}
	resp, err := client.SayHello(ctx, &in, nil)

	logrus.Infof("response:%v", resp.Message)
}

type yamuxSessionStream struct {
	net.Conn
	session *yamux.Session
}

func (y *yamuxSessionStream) Close() error {
	waitCh := y.session.CloseChan()
	timeout := time.NewTimer(defaultCloseTimeout)

	if err := y.Conn.Close(); err != nil {
		return err
	}

	if err := y.session.Close(); err != nil {
		return err
	}

	// block until session is really closed
	select {
	case <-waitCh:
		timeout.Stop()
	case <-timeout.C:
		return fmt.Errorf("timeout waiting for session close")
	}

	return nil
}

// NewAgentClient creates a new agent gRPC client and handles both unix and vsock addresses.
//
// Supported sock address formats are:
//   - unix://<unix socket path>
//   - vsock://<cid>:<port>
//   - <unix socket path>
func NewAgentClient(ctx context.Context, sock string, enableYamux bool) (*AgentClient, error) {
	grpcAddr, parsedAddr, err := parse(sock)
	if err != nil {
		return nil, err
	}
	dialOpts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	dialOpts = append(dialOpts, grpc.WithDialer(agentDialer(parsedAddr, enableYamux)))

	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &AgentClient{
		GreeterClient: pb.NewGreeterClient(conn),
		conn:          conn,
	}, nil
}


type dialer func(string, time.Duration) (net.Conn, error)

func agentDialer(addr *url.URL, enableYamux bool) dialer {
	var d dialer = unixDialer

	// yamux dialer
	return func(sock string, timeout time.Duration) (net.Conn, error) {
		conn, err := d(sock, timeout)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err != nil {
				conn.Close()
			}
		}()

		var session *yamux.Session
		sessionConfig := yamux.DefaultConfig()
		// Disable keepAlive since we don't know how much time a container can be paused
		sessionConfig.EnableKeepAlive = false
		sessionConfig.ConnectionWriteTimeout = time.Second
		session, err = yamux.Client(conn, sessionConfig)
		if err != nil {
			return nil, err
		}


		var stream net.Conn
		stream, err = session.Open()
		if err != nil {
			return nil, err
		}

		y := &yamuxSessionStream{
			Conn:    stream.(net.Conn),
			session: session,
		}

		return y, nil
	}
}

func unixDialer(sock string, timeout time.Duration) (net.Conn, error) {
	if strings.HasPrefix(sock, "unix:") {
		sock = strings.Trim(sock, "unix:")
	}

	dialFunc := func() (net.Conn, error) {
		return net.DialTimeout("unix", sock, timeout)
	}

	timeoutErr := grpcStatus.Errorf(codes.DeadlineExceeded, "timed out connecting to unix socket %s", sock)
	return commonDialer(timeout, dialFunc, timeoutErr)
}

// This would bypass the grpc dialer backoff strategy and handle dial timeout
// internally. Because we do not have a large number of concurrent dialers,
// it is not reasonable to have such aggressive backoffs which would kill kata
// containers boot up speed. For more information, see
// https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md
func commonDialer(timeout time.Duration, dialFunc func() (net.Conn, error), timeoutErrMsg error) (net.Conn, error) {
	t := time.NewTimer(timeout)
	cancel := make(chan bool)
	ch := make(chan net.Conn)
	go func() {
		for {
			select {
			case <-cancel:
				// canceled or channel closed
				return
			default:
			}

			conn, err := dialFunc()
			if err == nil {
				// Send conn back iff timer is not fired
				// Otherwise there might be no one left reading it
				if t.Stop() {
					ch <- conn
				} else {
					conn.Close()
				}
				return
			}
		}
	}()

	var conn net.Conn
	var ok bool
	select {
	case conn, ok = <-ch:
		if !ok {
			return nil, timeoutErrMsg
		}
	case <-t.C:
		cancel <- true
		return nil, timeoutErrMsg
	}

	return conn, nil
}


func parse(sock string) (string, *url.URL, error) {
	addr, err := url.Parse(sock)
	if err != nil {
		return "", nil, err
	}

	var grpcAddr string
	// validate more
	switch addr.Scheme {
	case unixSocketScheme:
		fallthrough
	case "":
		if (addr.Host == "" && addr.Path == "") || addr.Port() != "" {
			return "", nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid unix scheme: %s", sock)
		}
		if addr.Host == "" {
			grpcAddr = unixSocketScheme + ":///" + addr.Path
		} else {
			grpcAddr = unixSocketScheme + ":///" + addr.Host + "/" + addr.Path
		}
	default:
		return "", nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid scheme: %s", sock)
	}

	return grpcAddr, addr, nil
}
