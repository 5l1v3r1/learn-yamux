package main
// 多路复用
import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	pb "github.com/jimmy-xu/learn-yamux/yamux_serial_grpc/protocols/grpc"
)

// Set to the context that should be used for tracing gRPC calls.
var (
	grpcContext context.Context
	channelCloseTimeout    = 5 * time.Second
)

type Server struct{}
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	//get client info
	p, ok := peer.FromContext(ctx)
	if !ok {
		logrus.Errorf("failed to get peer of client")
	}
	logrus.Printf("receive gRPC request: [%v] client:%v", in.Name, p.Addr.String())
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

func Recv(stream net.Conn, id int){
	for {
		buf := make([]byte, 4)
		n, err := stream.Read(buf)
		if err == nil{
			fmt.Printf("ID:%v , len:%v unixtime:%v buf:%v\n", id, n, time.Now().Unix(), string(buf))
		}else{
			fmt.Printf("ID:%v , unixtime:%v error:%v\n", id, time.Now().Unix(), err)
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
		logrus.WithError(err).Fatal("Failed to create agent gRPC listener")
	}


	var (
		grpcServer *grpc.Server
	 	serverOpts []grpc.ServerOption
		servErr error
	)

	serverOpts = append(serverOpts, grpc.UnaryInterceptor(makeUnaryInterceptor()))
	grpcServer = grpc.NewServer(serverOpts...)
	pb.RegisterGreeterServer(grpcServer, &Server{})


	// session is closed when Serve() returns
	logrus.Infof("grpc server serve on yamux session")
	servErr = grpcServer.Serve(session)
	if servErr != nil {
		logrus.WithError(servErr).Warn("agent grpc server quits")
	}

	//
	//id :=0
	//for {
	//	// 建立多个流通路
	//	logrus.Printf("session Accept:%v", id)
	//	stream, err := session.Accept()
	//	if err == nil {
	//		logrus.Printf("Recv:%v", id)
	//		go Recv(stream, id)
	//		id ++
	//	}else{
	//		logrus.Println("session over.")
	//		break
	//	}
	//}

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
			logrus.Infof("waitCh")
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


func makeUnaryInterceptor() grpc.UnaryServerInterceptor {
	logrus.Infof("return makeUnaryInterceptor()")
	return func(origCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logrus.Infof("start makeUnaryInterceptor()")
		var start time.Time
		var elapsed time.Duration
		var message proto.Message

		grpcCall := info.FullMethod

		// Just log call details
		message = req.(proto.Message)

		logrus.WithFields(logrus.Fields{
			"request": grpcCall,
			"req":     message.String()}).Debug("new request")
		start = time.Now()


		// Use the context which will provide the correct trace
		// ordering, *NOT* the context provided to the function
		// returned by this function.
		resp, err = handler(getGRPCContext(), req)


		// Just log call details
		elapsed = time.Since(start)
		message = resp.(proto.Message)

		logger := logrus.WithFields(logrus.Fields{
			"request":  info.FullMethod,
			"duration": elapsed.String(),
			"resp":     message.String()})
		logger.Debug("request end")


		return resp, err
	}
}

func getGRPCContext() context.Context {
	if grpcContext != nil {
		return grpcContext
	}

	logrus.Warn("Creating gRPC context as none found")

	return context.Background()
}
