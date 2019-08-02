package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
)

type Server struct {
	listenAddr string
	done       chan<- bool
}

type Client struct {
	server   *Server
	clientIP string
}

type Stream struct {
	client   *Client
	streamID string
}

func main() {
	listenAddr := "0.0.0.0:4444"
	fmt.Println("[main]Starting yamux demo")
	done := make(chan bool, 0)

	server := NewServer(listenAddr, done)
	go server.Serv(listenAddr, done)

	fmt.Println("[main]wait done here")
	<-done
	fmt.Println("[main]close done ok")
}

func NewServer(listenAddr string, done chan<- bool) Server {
	return Server{
		listenAddr: listenAddr,
		done:       done,
	}
}

func (s *Server) Serv(listenAddr string, done chan bool) error {
	// Accept a TCP connection
	listener, err := net.Listen("tcp", listenAddr)
	defer close(done)

	logrus.Infof("start listener accept on %v", listenAddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Infof(fmt.Sprintf(" listener accept error:%v", err))
			return err
		}
		c := Client{
			server:   s,
			clientIP: conn.RemoteAddr().String(),
		}
		go c.handleConnection(&conn)
	}
	return err
}

func (c *Client) handleConnection(conn *net.Conn) {
	c.logger().Infof("\n== new client connected, creating yamux server session now ==")
	session, err := yamux.Server(*conn, nil)
	if err != nil {
		log.Printf(fmt.Sprintf("yamux server failed, error:%v", err))
		return
	}

	for {
		c.logger().Debugf("accepting yamux session stream...")
		stream, err := session.AcceptStream()
		if err != nil && err != io.EOF && !strings.Contains(err.Error(), "connection reset by peer") {
			c.logger().Errorf(fmt.Sprintf(" stop accept stream, error:%v", err))
			break
		}
		if stream != nil {
			t := Stream{
				client:   c,
				streamID: fmt.Sprintf("%v", stream.StreamID()),
			}
			go t.handleStream(stream)
		}
		if err != nil {
			c.logger().Infof("session closed")
			break
		}
	}
}

func (c *Client) logger() *logrus.Entry {
	return logrus.WithField("clientID", c.clientIP)
}

func (t *Stream) handleStream(stream *yamux.Stream) {
	if stream == nil {
		return
	}
	t.streamID = fmt.Sprintf("%v", stream.StreamID())
	t.logger().Debug("===== new stream opened =====")
	recvBuf := make([]byte, 32*1024)
	resultBuf := []byte{}
	for {
		n, err := (*stream).Read(recvBuf)
		if err != nil && err != io.EOF {
			t.logger().Errorf(fmt.Sprintf("read stream error:%v", err))
			break
		}
		resultBuf = t.processMessage(resultBuf, recvBuf[:n])
		if err == io.EOF {
			t.logger().Infof(string(resultBuf))
			t.logger().Debugf("stream closed")
			break
		}
	}
}

func (t *Stream) processMessage(buf, message []byte) []byte {
	return append(buf, message...)
}

func (t *Stream) logger() *logrus.Entry {
	return logrus.WithField("clientID", t.client.clientIP).WithField("streamID", t.streamID)
}
