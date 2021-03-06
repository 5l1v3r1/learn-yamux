package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func main() {
	var (
		conn    net.Conn
		session *yamux.Session
		stream  net.Conn
		err     error
		done    chan bool
	)

	done = make(chan bool)

	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				done <- false
				logrus.Infof("interrupt, signal:%v [error:%v]", s.String(), err)
			default:
				logrus.Infof("other signal", s)
			}
		}
	}()

	sock := ""
	if len(os.Args) > 1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock"
	}
	logrus.Printf("connect to %v", sock)

	var defaultDialTimeout = 5 * time.Second

	conn, err = unixDialer(sock, defaultDialTimeout)
	if err != nil {
		logrus.Fatalf("unix dialer failed, err:%v", err)
	}
	defer conn.Close()
	logrus.Infof("unix dial ok")

	defer func() {
		logrus.Infof("exit with error:%v", err)
	}()

	sessionConfig := yamux.DefaultConfig()
	// Disable keepAlive since we don't know how much time a container can be paused
	sessionConfig.EnableKeepAlive = false
	sessionConfig.ConnectionWriteTimeout = time.Second
	session, err = yamux.Client(conn, sessionConfig)
	if err != nil {
		logrus.Fatalf("create yamux client failed, error:%v", err)
	}
	defer session.Close()

	logrus.Infof("create yamux client ok")

	// 建立应用流通道1
	stream, err = session.Open()
	if err != nil {
		logrus.Fatalf("open session failed, err:%$v", err)
	}
	defer stream.Close()

	go func() {
		for i := 0; i < 3; i++ {
			logrus.Infof("send ping_%v", i)
			stream.Write([]byte(fmt.Sprintf("ping_%v", i)))
			buf := make([]byte, 128)
			if n, err := stream.Read(buf); err != nil {
				logrus.Error("failed to receive response, error:%v", err)
			} else {
				logrus.Infof("receive response: [%v]", string(buf[:n]))
			}
			time.Sleep(1 * time.Second)
		}
		done <- true
	}()

	rlt := <-done
	if rlt {
		logrus.Infof("exit normal")
	} else {
		logrus.Warnf("exit abnormal")
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
