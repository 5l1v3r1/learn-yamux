package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"

	"github.com/jimmy-xu/learn-yamux/pkg/serial"
)

var channelCloseTimeout = 5 * time.Second

func Recv(stream net.Conn, id int) {
	logrus.Info("Recv")
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err == nil {
			msg := string(buf[:n])
			logrus.Infof("recv [%s]", msg)
			if _, err = stream.Write([]byte(fmt.Sprintf("pong_%v", strings.Split(msg, "_")[1]))); err != nil {
				logrus.Errorf("failed to send response, error:%v", err)
			}
			logrus.Infof("send response [pong]")
		} else {
			if err == io.EOF {
				logrus.Errorf("stream read EOF")
				break
			} else {
				logrus.Errorf("failed to read stream, error:%v", err)
				break
			}
		}
	}
	logrus.Infof("close stream")
	stream.Close()
}
func main() {

	com := ""
	if len(os.Args) > 1 {
		com = os.Args[1]
	}

	if com == "" {
		//com = "COM1"
		//com = `\\.\agent.channel.0`
		com = `\\.\Global\agent.channel.0`
	}
	logrus.Printf("connect to %v", com)

	logrus.Infof("start setup()")
	var (
		sCh = &serialChannel{}
		err error
	)
	sCh.serialPath = com

	for {
		err = sCh.setup()
		if err != nil {
			logrus.Fatalf("setup() failed, error:%v", err)
		}

		logrus.Infof("===== start listen() =====")
		session, err := sCh.listen()
		if err != nil {
			logrus.WithError(err).Fatal("Failed to create agent grpc listener")
		}

		id := 0
		logrus.Println("session Accept")
		for {
			// 建立多个流通路
			stream, err := session.Accept()
			if err == nil {
				id++
				go Recv(stream, id)
			} else {
				logrus.Println("session is over.")
				break
			}
		}
		logrus.Infof("session close")
		session.Close()

		sCh.teardown()

		fmt.Println("<wait for new session>")
	}

	logrus.Infof("start teardown()")
	err = sCh.teardown()
	if err != nil {
		logrus.WithError(err).Warn("agent grpc channel teardown failed")
	}
}

type serialChannel struct {
	serialPath string
	serialConn *serial.Port
	waitCh     <-chan struct{}
}

func (c *serialChannel) setup() error {
	// Open serial channel.
	com := &serial.Config{Name: c.serialPath}
	s, err := serial.OpenPort(com)
	if err != nil {
		logrus.Fatalf("failed to open serial port %v, error:%v", c.serialPath, err)
	}
	logrus.Infof("open serial port %v ok", c.serialPath)
	c.serialConn = s

	return err
}

func (c *serialChannel) listen() (net.Listener, error) {
	config := yamux.DefaultConfig()
	// yamux client runs on the proxy side, sometimes the client is
	// handling other requests and it's not able to response to the
	// ping sent by the server and the communication is closed. To
	// avoid any IO timeouts in the communication between agent and
	// proxy, keep alive should be disabled.
	config.EnableKeepAlive = false
	config.LogOutput = yamuxWriter{}

	// Initialize Yamux server.
	session, err := yamux.Server(c.serialConn, config)
	if err != nil {
		return nil, err
	}
	logrus.Infof("init yamux server over serialport:%v ok", c.serialPath)
	c.waitCh = session.CloseChan()

	return session, nil
}

func (c *serialChannel) teardown() error {
	// wait for the session to be fully shutdown first
	if c.waitCh != nil {
		t := time.NewTimer(channelCloseTimeout)
		select {
		case <-c.waitCh:
			t.Stop()
		case <-t.C:
			return fmt.Errorf("timeout waiting for yamux channel to close")
		}
	}
	return c.serialConn.Close()
}

// yamuxWriter is a type responsible for logging yamux messages to the agent
// log.
type yamuxWriter struct {
}

// Write implements the Writer interface for the yamuxWriter.
func (yw yamuxWriter) Write(bytes []byte) (int, error) {
	message := string(bytes)

	l := len(message)

	// yamux messages are all warnings and errors
	logrus.WithField("component", "yamux").Warn(message)

	return l, nil
}
