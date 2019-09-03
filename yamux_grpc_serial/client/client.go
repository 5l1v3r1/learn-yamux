package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"

	pb "github.com/jimmy-xu/learn-yamux/pkg/grpc/protos"
)

const (
	unixSocketScheme  = "unix"
	vsockSocketScheme = "vsock"
)

var (
	defaultDialTimeout  = 5 * time.Second
	defaultCloseTimeout = 5 * time.Second

	checkRequestTimeout = 5 * time.Second
)

func main() {

	sock := ""
	if len(os.Args) > 1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/run/vc/vm/32272b12f09a0d91eae36c93773d8f8be17762003ccb9fac4446b3927742d242/kata.sock"
	}
	//logrus.Debugf("connect to %v", sock)

	agent := &kataAgent{}
	agent.state.URL = sock
	agent.ctx = context.Background()
	agent.keepConn = true

	for i := 0; i < 3; i++ {
		logrus.Println()
		logrus.Infof("===== send request %v =====", i)
		req := &pb.HelloRequest{Name: fmt.Sprintf("world %v", i)}
		resultingInterfaces, err := agent.sendReq(req)
		if err != nil {
			logrus.Warnf("failed to send grpc request, error:%v", err)
			continue
		}
		logrus.Infof("sent request: %v", req)
		resultInterfaces, ok := resultingInterfaces.(*pb.HelloResponse)
		if !ok {
			logrus.Fatalf("failed to get result, ok:%v", ok)
		}
		logrus.Infof("receive response: %v", resultInterfaces.Message)
	}

	if err := agent.disconnect(); err != nil {
		logrus.Fatalf("exit with error:%v", err)
	}
	logrus.Println()
	logrus.Infof("done")
	time.Sleep(1 * time.Second)
}

// AgentClient is an agent gRPC client connection wrapper for pb.AgentServiceClient
type AgentClient struct {
	pb.GreeterClient
	conn *grpc.ClientConn
}

func NewAgentClient(ctx context.Context, sock string, enableYamux bool) (*AgentClient, error) {
	grpcAddr, parsedAddr, err := parse(sock)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("grpcAddr:%v parsedAddr:%v", grpcAddr, parsedAddr)

	dialOpts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	dialOpts = append(dialOpts, grpc.WithDialer(agentDialer(parsedAddr, enableYamux)))

	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	logrus.Debugf("before grpc.DialContext: grpcAddr:%v", grpcAddr)
	conn, err := grpc.DialContext(ctx, grpcAddr, dialOpts...)
	if err != nil {
		return nil, err
	}

	logrus.Infof("grpc connection ok")
	return &AgentClient{
		GreeterClient: pb.NewGreeterClient(conn),
		conn:          conn,
	}, nil
}

// Close an existing connection to the agent gRPC server.
func (c *AgentClient) Close() error {
	logrus.Debugf("agent close")
	return c.conn.Close()
}

func unixDialer(sock string, timeout time.Duration) (net.Conn, error) {
	logrus.Debugf("start unixDialer()")
	if strings.HasPrefix(sock, "unix:") {
		sock = strings.Trim(sock, "unix:")
	}

	dialFunc := func() (net.Conn, error) {
		logrus.Debugf("start net.DialTimeout sock:%v timeout:%v", sock, timeout)
		return net.DialTimeout("unix", sock, timeout)
	}

	timeoutErr := grpcStatus.Errorf(codes.DeadlineExceeded, "timed out connecting to unix socket %s", sock)
	return commonDialer(timeout, dialFunc, timeoutErr)
}

func commonDialer(timeout time.Duration, dialFunc func() (net.Conn, error), timeoutErrMsg error) (net.Conn, error) {
	t := time.NewTimer(timeout)
	cancel := make(chan bool)
	ch := make(chan net.Conn)
	go func() {
		for {
			select {
			case <-cancel:
				logrus.Info("canceled or channel closed")
				return
			default:
				logrus.Debugf("waiting...")
			}

			conn, err := dialFunc()
			if err == nil {
				// Send conn back if timer is not fired
				// Otherwise there might be no one left reading it
				if t.Stop() {
					logrus.Debugf("commonDialer conn ok")
					ch <- conn
				} else {
					logrus.Info("commonDialer conn close")
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
		logrus.Infof("unix socket connected")
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
	case vsockSocketScheme:
		if addr.Hostname() == "" || addr.Port() == "" || addr.Path != "" {
			return "", nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid vsock scheme: %s", sock)
		}
		if _, err := strconv.ParseUint(addr.Hostname(), 10, 32); err != nil {
			return "", nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid vsock cid: %s", sock)
		}
		if _, err := strconv.ParseUint(addr.Port(), 10, 32); err != nil {
			return "", nil, grpcStatus.Errorf(codes.InvalidArgument, "Invalid vsock port: %s", sock)
		}
		grpcAddr = vsockSocketScheme + ":" + addr.Host
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

type dialer func(string, time.Duration) (net.Conn, error)

func agentDialer(addr *url.URL, enableYamux bool) dialer {
	logrus.Debugf("start agentDialer()")
	var d dialer
	switch addr.Scheme {
	case unixSocketScheme:
		fallthrough
	default:
		d = unixDialer
	}

	if !enableYamux {
		return d
	}

	logrus.Debugf("return yamux dialer")

	// yamux dialer
	return func(sock string, timeout time.Duration) (net.Conn, error) {
		logrus.Debugf("start yamux dialer, sock:%v timeout:%v", sock, timeout)
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
		logrus.Infof("create yamux session")
		session, err = yamux.Client(conn, sessionConfig)
		if err != nil {
			return nil, err
		}

		var stream net.Conn
		logrus.Infof("create yamux stream")
		stream, err = session.OpenStream()
		if err != nil {
			logrus.Errorf("failed to create yamux stream, error:%v", err)
			return nil, err
		}

		logrus.Debugf("yamux create stream ok")
		y := &yamuxSessionStream{
			Conn:    stream.(net.Conn),
			session: session,
		}

		return y, nil
	}
}

type yamuxSessionStream struct {
	net.Conn
	session *yamux.Session
}

func (y *yamuxSessionStream) Close() error {
	logrus.Infof("yamux session close")
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

type KataAgentState struct {
	ProxyPid int
	URL      string
}

type kataAgent struct {
	// lock protects the client pointer
	sync.Mutex
	client *AgentClient

	reqHandlers  map[string]reqFunc
	state        KataAgentState
	keepConn     bool
	proxyBuiltIn bool

	vmSocket interface{}
	ctx      context.Context
}

func (k *kataAgent) connect() error {
	logrus.Debugf("start connect")

	// This is for the first connection only, to prevent race
	k.Lock()
	defer k.Unlock()

	// lockless quick pass
	if k.client != nil {
		logrus.Infof("use exist kata agent client")
		return nil
	} else {
		logrus.Infof("client is nil, create a new one")
	}

	kataURL := k.state.URL
	logrus.WithField("url", kataURL).Info("New client")
	if k.ctx == nil {
		logrus.WithField("type", "bug").Error("trace called before context set")
		k.ctx = context.Background()
	}

	logrus.Debugf("NewAgentClient, kataURL: %v", kataURL)
	client, err := NewAgentClient(k.ctx, kataURL, true)
	if err != nil {
		return err
	}

	k.installReqFunc(client)
	k.client = client

	return nil
}

type reqFunc func(context.Context, interface{}, ...grpc.CallOption) (interface{}, error)

func (k *kataAgent) installReqFunc(c *AgentClient) {
	logrus.Debugf("start installReqFunc")

	k.reqHandlers = make(map[string]reqFunc)
	k.reqHandlers["grpc.SayHello"] = func(ctx context.Context, req interface{}, opts ...grpc.CallOption) (interface{}, error) {
		return k.client.SayHello(ctx, req.(*pb.HelloRequest), opts...)
	}
}

func (k *kataAgent) sendReq(request interface{}) (interface{}, error) {
	logrus.Debugf("start sendReq")
	if err := k.connect(); err != nil {
		return nil, err
	}
	logrus.Debugf("connect ok")
	if !k.keepConn {
		defer k.disconnect()
	}

	logrus.Debugf("reqHandlers:%v", k.reqHandlers)
	logrus.Debugf("request:%v", request)

	msgName := proto.MessageName(request.(proto.Message))
	if msgName == "" {
		msgName = "grpc.SayHello"
	}
	logrus.Debugf("get handler for %v", msgName)
	handler := k.reqHandlers[msgName]
	if msgName == "" || handler == nil {
		return nil, errors.New("Invalid request type")
	}
	message := request.(proto.Message)
	logrus.WithField("name", msgName).WithField("req", message.String()).Debug("sending request")

	logrus.Debugf("call handler")
	return handler(k.ctx, request)
}

func (k *kataAgent) disconnect() error {
	k.Lock()
	defer k.Unlock()

	if k.client == nil {
		return nil
	}

	if err := k.client.Close(); err != nil && grpcStatus.Convert(err).Code() != codes.Canceled {
		return err
	}

	k.client = nil
	k.reqHandlers = nil

	return nil
}
