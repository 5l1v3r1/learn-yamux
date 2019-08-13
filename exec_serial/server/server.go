package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/jimmy-xu/learn-yamux/pkg/serial"
	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Starting exec_serial demo - server")

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
	fmt.Printf("[read] open serial port %v ok\n", kataPort)

	////////////////////////////////////////////////
	//console.sock
	consolePort := `\\.\console0`

	console := &serial.Config{Name: consolePort}
	c, err := serial.OpenPort(console)
	if err != nil {
		logrus.Fatalf("failed to open console port %v, error:%v", consolePort, err)
	}
	fmt.Printf("[write] open serial port %v ok\n", consolePort)

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
				//splitting a string by space, except inside quotes
				re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
				arr := re.FindAllString(string(buf[:n-1]), -1)
				var (
					cmd         *exec.Cmd
					interactive = true
				)
				if interactive {
					logrus.Infof("interactive mode")
					opt := []string{"/k"}
					for _, v := range arr {
						opt = append(opt, strings.Trim(v, `"`))
					}
					cmd = exec.Command("cmd", opt...)
					logrus.Infof("execute cmd:%v", cmd.Args)
					cmd.Stdin = s
					cmd.Stderr = c
					cmd.Stdout = c
					cmd.Run()
					logrus.Infof("finish execute interactive cmd:%v", cmd.Args)
				} else {
					logrus.Infof("non-interactive mode")
					opt := []string{"/c"}
					for _, v := range arr {
						opt = append(opt, strings.Trim(v, `"`))
					}
					cmd = exec.Command("cmd", opt...)
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
					c.Write([]byte(`C:\Users\admin>`))
					logrus.Infof("finish execute non-interactive cmd:%v", cmd.Args)
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
