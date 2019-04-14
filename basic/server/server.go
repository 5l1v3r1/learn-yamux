package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

func main() {
	localAddr := ":4444"
	fmt.Printf("Starting yamux demo - server on %v\n", localAddr)
	done := make(chan bool, 0)
	go server(localAddr, done)

	fmt.Println("wait done here")
	<-done

	time.Sleep(time.Second * 20)
}

func server(localAddr string, done chan bool) error {
	// Accept a TCP connection
	listener, err := net.Listen("tcp", localAddr)

	close(done)
	fmt.Println("close done ok")

	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	// Setup server side of yamux
	log.Println("creating server session")
	session, err := yamux.Server(conn, nil)
	if err != nil {
		return err
	}

	// Accept a stream
	log.Println("accepting stream")
	stream, err := session.Accept()
	if err != nil {
		return err
	}

	// Listen for a message
	buf := make([]byte, 255)
	n, err := stream.Read(buf)

	fmt.Printf("n:%v err:%v buf = %+v", n, err, string(buf))

	return err
}

