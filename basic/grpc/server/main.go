package main

import (
	"net"

	pb "github.com/jimmy-xu/learn-yamux/pkg/grpc/protos"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = ":50051"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	//get client info
	p, ok := peer.FromContext(ctx)
	if !ok {
		logrus.Errorf("failed to get peer of client")
	}
	logrus.Printf("[%v]client connected, received request: [%v], sent response [hello %v]", p.Addr.String(), in.Name, in.Name)
	return &pb.HelloResponse{Message: "hello " + in.Name}, nil
}
func main() {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logrus.Printf("failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &Server{})
	reflection.Register(grpcServer)

	logrus.Printf("Start gRPC Server on port :%v", grpcPort)
	grpcServer.Serve(listen)
}
