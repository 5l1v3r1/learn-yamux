package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting serial demo - client")

	sock := ""
	if len(os.Args) > 1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock"
	}
	logrus.Printf("connect to %v", sock)

	var defaultDialTimeout = 3 * time.Second

	//use unix sock file
	conn, err := unixDialer(sock, defaultDialTimeout)
	if err != nil {
		logrus.Fatalf("unix dialer failed, err:%v", err)
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()
	logrus.Infof("unix dial ok")

	logrus.Info("===== send message to server =====")
	for idx, item := range []string{"hello", "world", "ping"} {
		if _, err := conn.Write([]byte(item)); err != nil {
			logrus.Error("[%v]failed to send %v, error:%v", idx, item, err)
		} else {
			logrus.Infof("[%v]sent '%v'", idx, item)
			time.Sleep(10 * time.Millisecond)
		}
	}

	logrus.Info("===== read message from server =====")
	buf := make([]byte, 1024)
	if n, err := conn.Read(buf); err != nil {
		logrus.Error("failed to read from server, error:%v", err)
	} else {
		logrus.Infof("received :%s", buf[:n])
	}
}

func unixDialer(sock string, timeout time.Duration) (net.Conn, error) {
	if strings.HasPrefix(sock, "unix:") {
		sock = strings.Trim(sock, "unix:")
	}

	dialFunc := func() (net.Conn, error) {
		return net.DialTimeout("unix", sock, timeout)
	}

	return commonDialer(timeout, dialFunc, fmt.Errorf("timed out connecting to unix socket %s", sock))
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
