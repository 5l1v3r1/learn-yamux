package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting yamux demo - client")

	address := ""
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	if address == "" {
		address = "127.0.0.1:4444"
	}
	logrus.Printf("connect to %v", address)

	if err := client(address); err != nil {
		logrus.Fatal(err)
	}

}

func client(serverAddr string) error {
	// Get a TCP connection
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}

	// Setup client side of yamux
	logrus.Infof("creating client session")
	session, err := yamux.Client(conn, nil)
	if err != nil {
		return err
	}
	logrus.Infof("create yamux session ok")
	defer session.Close()

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		//goroutine for each stream
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Open a new stream
			//stream, err := session.Open()
			stream, err := session.OpenStream()
			if err != nil {
				log.Printf("failed to open stream, error:%v", err)
				return
			}
			defer stream.Close()
			logrus.Infof("opening stream %v", stream.StreamID())

			// Stream implements net.Conn
			for idx, item := range []string{"hello", " ", "world"} {
				sendMessage(stream, idx, item)
			}

		}()
	}
	logrus.Infof("waiting...")
	wg.Wait()
	logrus.Info("done")
	return nil
}

func sendMessage(stream *yamux.Stream, idx int, message string) error {
	_, err := stream.Write([]byte(fmt.Sprintf("[%v]%v", idx, message)))
	if err != nil {
		logrus.Warnf("write message error:%v", err)
	}
	return nil
}
