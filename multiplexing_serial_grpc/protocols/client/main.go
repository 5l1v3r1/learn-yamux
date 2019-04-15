package main

import (
	"context"
	"os"
	"time"

	pb "github.com/jimmy-xu/learn-yamux/multiplexing_serial_grpc/protocols/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

func main() {

	address := ""
	if len(os.Args)>1 {
		address = os.Args[1]
	}

	if address == "" {
		address = "127.0.0.1:50051"
	}
	logrus.Printf("connect to %v", address)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		logrus.Fatalf("could not greet: %v", err)
	}
	logrus.Printf("Receiver gRPC response: %s", r.Message)
}