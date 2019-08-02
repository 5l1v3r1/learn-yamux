package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jimmy-xu/learn-yamux/pkg/serial"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting serial demo - server")

	serialPort := ""
	if len(os.Args) > 1 {
		serialPort = os.Args[1]
	}
	if serialPort == "" {
		serialPort = `\\.\agent.channel.0`
	}

	com := &serial.Config{Name: serialPort}
	s, err := serial.OpenPort(com)
	if err != nil {
		logrus.Fatalf("failed to open serial port %v, error:%v", serialPort, err)
	}
	fmt.Printf("open serial port %v ok", serialPort)

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 1024)
		var readCount int

		logrus.Info("===== begin to receive message from client ======")
		for {
			n, err := s.Read(buf)
			if err != nil {
				if strings.Contains(err.Error(), "Insufficient system resources exist to complete the requested service") {
					time.Sleep(1 * time.Second)
					continue
				} else {
					logrus.Fatalf("failed to read serial port, error:%v", err)
					break
				}
			}
			readCount++
			logrus.Infof("[%v]received: %s", readCount, buf[:n])
			if string(buf[:n]) == "ping" {
				if _, err := s.Write([]byte("pong")); err != nil {
					logrus.Error("failed to send pong to client, error:%v", err)
				} else {
					logrus.Info("sent: pong")
				}
			}

			select {
			case <-ch:
				ch <- readCount
				close(ch)
			default:
			}
		}
	}()

	fmt.Println("wait...")
	<-ch
	fmt.Println("done")
}
