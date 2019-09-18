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
	defaultDialTimeout = 2 * time.Second
	defer func() {
		defaultDialTimeout = dialTimeout
	}()
	var (
		cli *AgentClient
		err error
	)
	for i := 0; i < 6; i++ {
		cli, err = NewAgentClient(context.Background(), sock, enableYamux)
		if err == nil {
			break
		}
		logrus.Warnf("%v: failed to create new agent client, retry", i)
		continue
		time.Sleep(1 * time.Second)
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
		logrus.Fatalf("failed to get guest details: %s", err)
	}

	err = getNetworkConfig(cli)
	if err != nil {
		logrus.Fatalf("failed to get network config: %s", err)
	}

	err = getUsers(cli)
	if err != nil {
		logrus.Fatalf("failed to get users: %s", err)
	}

	err = setHostname(cli)
	if err != nil {
		logrus.Fatalf("failed to set hostname: %s", err)
	}
	err = getHostname(cli)
	if err != nil {
		logrus.Fatalf("failed to get hostname: %s", err)
	}

	err = setKMS(cli)
	if err != nil {
		logrus.Fatalf("failed to set kms server: %s", err)
	}
	err = getKMS(cli)
	if err != nil {
		logrus.Fatalf("failed to get kms server: %s", err)
	}

	err = setUserPassword(cli)
	if err != nil {
		logrus.Fatalf("failed to set user password: %s", err)
	}
}

func checkVersion(cli *AgentClient) error {
	logrus.Infof("---------- [request] health.Version() ----------")
	resp, err := cli.Version(context.Background(), &pb.CheckRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
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
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func getGuestDetails(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetGuestDetails() ----------")
	resp, err := cli.GetGuestDetails(context.Background(), &pb.GuestDetailsRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func getNetworkConfig(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetNetworkConfig() ----------")
	resp, err := cli.GetNetworkConfig(context.Background(), &pb.GetNetworkConfigRequest{MacAddress: "E6-29-C9-D2-00-0F"})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func getUsers(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetUsers() ----------")
	resp, err := cli.GetUsers(context.Background(), &pb.GetUsersRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func getHostname(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetHostname() ----------")
	resp, err := cli.GetHostname(context.Background(), &pb.GetHostnameRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func getKMS(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.GetKMS() ----------")
	resp, err := cli.GetKMS(context.Background(), &pb.GetKMSRequest{})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func setUserPassword(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.SetUserPassword() ----------")
	resp, err := cli.SetUserPassword(context.Background(), &pb.SetUserPasswordRequest{Username: "admin", Password: "Test123!@#"})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func setKMS(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.SetKMS() ----------")
	resp, err := cli.SetKMS(context.Background(), &pb.SetKMSRequest{Server: "kms.alibaba-inc.com"})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func setHostname(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.SetHostname() ----------")
	resp, err := cli.SetHostname(context.Background(), &pb.SetHostnameRequest{Hostname: "jimmy-win7"})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}
