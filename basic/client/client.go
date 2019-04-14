package main

import (
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/yamux"
)

func main() {
	fmt.Println("Starting yamux demo - client")

	localAddr := "172.16.87.100:4444"

	if err := client(localAddr); err != nil {
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
