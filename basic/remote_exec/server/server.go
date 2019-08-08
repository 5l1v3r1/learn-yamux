package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jimmy-xu/learn-yamux/pkg/serial"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting remote_exec demo - server")

	////////////////////////////////////////////////
	//kata.sock
	kataPort := ""
	if len(os.Args) > 1 {
		kataPort = os.Args[1]
	}
	if kataPort == "" {
		kataPort = `\\.\agent.channel.0`
	}
	com := &serial.Config{Name: kataPort}
	s, err := serial.OpenPort(com)
	if err != nil {
		logrus.Fatalf("failed to open serial port %v, error:%v", kataPort, err)
	}
	fmt.Printf("open serial port %v ok\n", kataPort)

	////////////////////////////////////////////////
	//console.sock
	consolePort := `\\.\console0`

	console := &serial.Config{Name: consolePort}
	c, err := serial.OpenPort(console)
	if err != nil {
		logrus.Fatalf("failed to open console port %v, error:%v", consolePort, err)
	}
	fmt.Printf("open console port %v ok\n", consolePort)

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
			} else {
				cmd := exec.Command("cmd", "/c", string(buf[:n]))
				logrus.Infof("execute cmd:%v", cmd.Args)

				stdout, _ := cmd.StdoutPipe()
				stderr, _ := cmd.StderrPipe()
				cmd.Start()
				oRlt := make([]byte, 1024)
				eRlt := make([]byte, 1024)
				for {
					oN, oErr := stdout.Read(oRlt)
					if oN > 0 {
						c.Write([]byte(fmt.Sprintf("%s", oRlt[:oN])))
					}
					if oErr != nil {
						eN, eErr := stderr.Read(eRlt)
						if oErr == io.EOF && eErr == io.EOF {
							break
						}
						if eN > 0 {
							c.Write([]byte(fmt.Sprintf("%s", eRlt[:eN])))
						}
					}
				}
				logrus.Infof("finish execute cmd:%v", cmd.Args)
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
