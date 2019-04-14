package main
// 多路复用
import (
	"fmt"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"net"
	"time"
)

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
	// 建立底层复用连接
	tcpaddr, _ := net.ResolveTCPAddr("tcp4", ":8980")
	tcplisten, _ := net.ListenTCP("tcp", tcpaddr)

	logrus.Println("TCP listen Accept")
	conn, _ := tcplisten.Accept()

	fmt.Println("yamux server on TCP")
	session, _ := yamux.Server(conn, nil)

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

}