package main
// 多路复用
import (
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"time"
)

var channelCloseTimeout    = 5 * time.Second

func Recv(stream net.Conn, id int){
	for {
		buf := make([]byte, 4)
		n, err := stream.Read(buf)
		if err == nil{
			fmt.Println("ID:", id, ", len:", n, time.Now().Unix(), string(buf))
		}else{
			fmt.Println(time.Now().Unix(), err)
			return
		}
	}
}
func main()  {

	com := ""
	if len(os.Args)>1 {
		com = os.Args[1]
	}

	if com == "" {
		com = "COM1"
	}
	logrus.Printf("connect to %v", com)


	logrus.Infof("start setup()")
	var sCh = &serialChannel{}
	sCh.serialPath = com
	err := sCh.setup()
	if err != nil {
		logrus.Fatalf("setup() failed, error:%v", err)
	}

	logrus.Infof("start listen()")
	session, err := sCh.listen()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create agent grpc listener")
	}

	id :=0
	for {
		// 建立多个流通路
		logrus.Println("session Accept")
		stream, err := session.Accept()
		if err == nil {
			logrus.Println("Recv")
			id ++
			go Recv(stream, id)
		}else{
			logrus.Println("session over.")
			return
		}
	}

	logrus.Infof("start teardown()")
	err = sCh.teardown()
	if err != nil {
		logrus.WithError(err).Warn("agent grpc channel teardown failed")
	}


}


type serialChannel struct {
	serialPath string
	serialConn *os.File
	waitCh     <-chan struct{}
}

func (c *serialChannel) setup() error {
	// Open serial channel.
	file, err := os.OpenFile(c.serialPath, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return err
	}

	c.serialConn = file

	return nil
}


func (c *serialChannel) listen() (net.Listener, error) {
	config := yamux.DefaultConfig()
	// yamux client runs on the proxy side, sometimes the client is
	// handling other requests and it's not able to response to the
	// ping sent by the server and the communication is closed. To
	// avoid any IO timeouts in the communication between agent and
	// proxy, keep alive should be disabled.
	config.EnableKeepAlive = false
	config.LogOutput = yamuxWriter{}

	// Initialize Yamux server.
	session, err := yamux.Server(c.serialConn, config)
	if err != nil {
		return nil, err
	}
	c.waitCh = session.CloseChan()

	return session, nil
}

func (c *serialChannel) teardown() error {
	// wait for the session to be fully shutdown first
	if c.waitCh != nil {
		t := time.NewTimer(channelCloseTimeout)
		select {
		case <-c.waitCh:
			t.Stop()
		case <-t.C:
			return fmt.Errorf("timeout waiting for yamux channel to close")
		}
	}
	return c.serialConn.Close()
}



// yamuxWriter is a type responsible for logging yamux messages to the agent
// log.
type yamuxWriter struct {
}

// Write implements the Writer interface for the yamuxWriter.
func (yw yamuxWriter) Write(bytes []byte) (int, error) {
	message := string(bytes)

	l := len(message)

	// yamux messages are all warnings and errors
	logrus.WithField("component", "yamux").Warn(message)

	return l, nil
}