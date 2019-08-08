package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
	GB2312  = Charset("GB2312")
)

func main() {
	fmt.Println("Starting remote_exec demo - client")

	var defaultDialTimeout = 3 * time.Second

	////////////////////////////////////////////////////////////
	//use unix kataSock file
	kataSock := "/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/kata.sock"
	kataConn, err := unixDialer(kataSock, defaultDialTimeout)
	if err != nil {
		logrus.Fatalf("unix dialer failed, err:%v", err)
	}
	logrus.Infof("unix dial %v ok", kataSock)

	////////////////////////////////////////////////////////////
	consoleSock := "/run/vc/vm/1cd65c2aefcb65ee2a2139373f4e041f35074b2d5a0f0c3f274ec2e9cdc18694/console.sock"
	consoleConn, err := unixDialer(consoleSock, defaultDialTimeout)
	if err != nil {
		logrus.Fatalf("unix dialer failed, err:%v", err)
	}
	logrus.Infof("unix dial %v ok", consoleSock)

	////////////////////////////////////////////////////////////
	defer func() {
		if err = kataConn.Close(); err != nil {
			logrus.Fatalf("failed to close kata sock, error:%v", err)
		}
		logrus.Infof("%v closed", kataSock)
		if err = consoleConn.Close(); err != nil {
			logrus.Fatalf("failed to close console sock, error:%v", err)
		}
		logrus.Infof("%v closed", consoleSock)
	}()

	//read
	fmt.Printf(`c:\>`)
	go func() {
		for {
			rlt := make([]byte, 1024)
			if n, err := consoleConn.Read(rlt); err != nil {
				if err == io.EOF {
					fmt.Printf("[%v] %s", err, convertByte2String(rlt[:n], GB18030))
				}
				fmt.Printf(`c:\>`)
			} else {
				fmt.Printf("%s", convertByte2String(rlt[:n], GB18030))
			}
		}
	}()

	//write
	for {
		bio := bufio.NewReader(os.Stdin)
		buf, _, _ := bio.ReadLine()
		buf = append(buf, '\n')
		kataConn.Write(buf)
	}
}

func unixDialer(sock string, timeout time.Duration) (net.Conn, error) {
	if strings.HasPrefix(sock, "unix:") {
		sock = strings.Trim(sock, "unix:")
	}

	dialFunc := func() (net.Conn, error) {
		return net.DialTimeout("unix", sock, timeout)
	}

	return commonDialer(timeout, dialFunc, fmt.Errorf("timed out connecting to unix socket %s", sock))
}

func commonDialer(timeout time.Duration, dialFunc func() (net.Conn, error), timeoutErrMsg error) (net.Conn, error) {
	t := time.NewTimer(timeout)
	cancel := make(chan bool)
	ch := make(chan net.Conn)
	go func() {
		for {
			select {
			case <-cancel:
				// canceled or channel closed
				return
			default:
			}

			conn, err := dialFunc()
			if err == nil {
				// Send conn back iff timer is not fired
				// Otherwise there might be no one left reading it
				if t.Stop() {
					ch <- conn
				} else {
					conn.Close()
				}
				return
			}
		}
	}()

	var conn net.Conn
	var ok bool
	select {
	case conn, ok = <-ch:
		if !ok {
			return nil, timeoutErrMsg
		}
	case <-t.C:
		cancel <- true
		return nil, timeoutErrMsg
	}

	return conn, nil
}

func convertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case GBK:
		var decodeBytes, _ = simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case GB2312:
		var decodeBytes, _ = simplifiedchinese.HZGB2312.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}
