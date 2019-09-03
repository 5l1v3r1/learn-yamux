package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/jimmy-xu/learn-yamux/pkg/grpc/protos"
)

func main() {
	sock := ""
	if len(os.Args) > 1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/run/vc/vm/32272b12f09a0d91eae36c93773d8f8be17762003ccb9fac4446b3927742d242/kata.sock"
	} else {
		sock = fmt.Sprintf("/run/vc/vm/%s/kata.sock", sock)
	}

	testAgentClient(sock, true)
}

func testAgentClient(sock string, enableYamux bool) {
	dialTimeout := defaultDialTimeout
	defaultDialTimeout = 5 * time.Second
	defer func() {
		defaultDialTimeout = dialTimeout
	}()
	cli, err := NewAgentClient(context.Background(), sock, enableYamux)
	if err != nil {
		logrus.Fatalf("Failed to create new agent client: %s", err)
	}

	logrus.Infof("connect to kata-agent server ok")
	defer cli.Close()

	err = checkVersion(cli)
	if err != nil {
		logrus.Fatalf("failed checking grpc server version: %s", err)
	}
}

func checkVersion(cli *AgentClient) error {
	resp, err := cli.GetGuestDetails(context.Background(), &pb.GuestDetailsRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("version:%v", resp.AgentDetails.Version)
	return nil
}


