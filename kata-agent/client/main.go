package main

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	defaultDialTimeout = 2 * time.Second
	defer func() {
		defaultDialTimeout = dialTimeout
	}()
	var (
		cli *AgentClient
		err error
	)
	for i := 0; i<6; i++ {
		cli, err = NewAgentClient(context.Background(), sock, enableYamux)
		if err == nil {
			break
		}
		logrus.Warnf("%v: failed to create new agent client, retry", i)
		continue
		time.Sleep(1*time.Second)
	}
	if err != nil {
		logrus.Fatalf("Failed to connect to agent client: %s", err)
	}

	logrus.Infof("connect to kata-agent server via %v ok", sock)
	defer cli.Close()

	err = checkHealth(cli)
	if err != nil {
		logrus.Fatalf("failed checking grpc server health: %s", err)
	}

	err = checkVersion(cli)
	if err != nil {
		logrus.Fatalf("failed checking grpc server version: %s", err)
	}

	err = getGuestDetails(cli)
	if err != nil {
		logrus.Fatalf("failed get guest details: %s", err)
	}
}

func checkVersion(cli *AgentClient) error {
	logrus.Infof("---------- [request] health.Version() ----------")
	resp, err := cli.Version(context.Background(), &pb.CheckRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", strings.Join(strings.Split(resp.String()," "),"\n"))
	return nil
}

func checkHealth(cli *AgentClient) error {
	logrus.Infof("---------- [request] health.Check() ----------")
	resp, err := cli.Check(context.Background(), &pb.CheckRequest{})
	if err != nil {
		return err
	}
	if resp.Status != pb.HealthCheckResponse_SERVING {
		return fmt.Errorf("unexpected health status: %s", resp.Status)
	}
	logrus.Infof("[response]:\n%s", strings.Join(strings.Split(resp.String()," "),"\n"))
	return nil
}

func getGuestDetails(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetGuestDetails() ----------")
	resp, err := cli.GetGuestDetails(context.Background(), &pb.GuestDetailsRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", strings.Join(strings.Split(resp.String()," "),"\n"))
	return nil
}
