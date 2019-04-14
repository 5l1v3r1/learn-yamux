package main

import (
	"net"
	"os"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
)

func main()  {

	address := ""
	if len(os.Args)>1 {
		address = os.Args[1]
	}

	if address == "" {
		address = "127.0.0.1:8980"
	}
	logrus.Printf("connect to %v", address)


	// 建立底层复用通道
	conn, _ := net.Dial("tcp", address)
	session, _ := yamux.Client(conn, nil)

	// 建立应用流通道1
	stream1, _ := session.Open()
	stream1.Write([]byte("ping" ))
	stream1.Write([]byte("pnng" ))
	time.Sleep(1 * time.Second)

	// 建立应用流通道2
	stream2, _ := session.Open()
	stream2.Write([]byte("pong"))
	time.Sleep(1 * time.Second)

	// 清理退出
	time.Sleep(5 * time.Second)

	stream1.Close()
	stream2.Close()

	session.Close()

	conn.Close()
}