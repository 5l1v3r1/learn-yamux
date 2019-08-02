package main

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
)

func main()  {

	sock := ""
	if len(os.Args)>1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/Users/xjimmy/Documents/antfin/serial/com1"
	}
	logrus.Printf("connect to %v", sock)

	var defaultDialTimeout = 15 * time.Second


	//use unix sock file
	conn, err := unixDialer(sock, defaultDialTimeout)
	if err != nil {
		logrus.Fatalf("unix dialer failed, err:%v", err)
	}

	logrus.Infof("unix dial ok")

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
		logrus.Fatalf("create yamux client failed, error:%v", err)
	}

	logrus.Infof("create yamux client ok")

	// 建立应用流通道1
	stream1, err := session.Open()
	if err != nil {
		logrus.Fatalf("open session failed, err:%$v", err)
	}

	logrus.Info("send ping")
	stream1.Write([]byte("ping" ))
	stream1.Write([]byte("pnng" ))
	time.Sleep(1 * time.Second)

	// 建立应用流通道2
	logrus.Info("send pong")
	stream2, _ := session.Open()
	stream2.Write([]byte("pong"))
	time.Sleep(1 * time.Second)

	// 清理退出
	time.Sleep(5 * time.Second)

	logrus.Info("close")
	stream1.Close()
	stream2.Close()

	session.Close()

	conn.Close()
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
