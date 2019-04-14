package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"

	"github.com/hashicorp/yamux"
)

func main() {
	fmt.Println("Starting yamux demo - client")

	address := ""
	if len(os.Args)>1 {
		address = os.Args[1]
	}

	if address == "" {
		address = "127.0.0.1:4444"
	}
	logrus.Printf("connect to %v", address)


	if err := client(address); err != nil {
		log.Println(err)
	}

}

func client(serverAddr string) error {
	// Get a TCP connection
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	// Setup client side of yamux
	log.Println("creating client session")
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}

	// Open a new stream
	log.Println("opening stream")
	stream, err := session.Open()
	if err != nil {
		return err
	}

	// Stream implements net.Conn
	_, err = stream.Write([]byte("hello world"))
	return err
}
