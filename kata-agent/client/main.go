package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	pb "github.com/jimmy-xu/learn-yamux/pkg/grpc/protos"
)

var (
	macAddr        = "9e:05:8e:c4:f1:70" // E6:29:C9:D2:00:0F e6:29:c9:d2:00:0f e6-29-c9-d2-00-0f
	hostname       = "jimmy-win7"
	ipAddr         = "172.19.0.245"
	subnet         = "255.255.255.0"
	defaultGateway = "172.19.0.1"
	dns            = []string{"172.16.87.1", "8.8.4.4"}
)

func main() {
	sock := ""
	if len(os.Args) > 1 {
		sock = os.Args[1]
	}

	if sock == "" {
		sock = "/run/vc/vm/ba9559adc4007b9f6f4788287f5086831cfe601af9225a3382899dcbeba17dca/kata.sock"
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

	err = setNetworkConfig(cli)
	if err != nil {
		logrus.Fatalf("failed to set network config: %s", err)
	}

	err = getNetworkConfig(cli)
	if err != nil {
		logrus.Fatalf("failed to get network config: %s", err)
	}

	err = getUsers(cli)
	if err != nil {
		logrus.Fatalf("failed to get users: %s", err)
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

	curHostname, err := getHostname(cli)
	if err != nil {
		logrus.Fatalf("failed to get hostname: %s", err)
	}

	if curHostname != hostname {
		logrus.Infof("current hostname is %v, rename to %v", curHostname, hostname)
		err = setHostname(cli)
		if err != nil {
			logrus.Fatalf("failed to set hostname: %s", err)
		}
	} else {
		logrus.Infof("hostname is %v, ignore", hostname)
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
	resp, err := cli.GetNetworkConfig(context.Background(), &pb.GetNetworkConfigRequest{MacAddress: macAddr})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}

func setNetworkConfig(cli *AgentClient) error {
	logrus.Infof("---------- [request] agent.SetNetworkConfig() ----------")
	resp, err := cli.SetNetworkConfig(context.Background(), &pb.SetNetworkConfigRequest{
		MacAddress: macAddr,
		Gateway:    defaultGateway,
		DnsServer:  dns,
		Addrs: []*pb.Addrs{
			{IpAddress: ipAddr, Subnet: subnet},
		},
	})
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

func getHostname(cli *AgentClient) (string, error) {
	logrus.Infof("---------- [request] agent.GetHostname() ----------")
	resp, err := cli.GetHostname(context.Background(), &pb.GetHostnameRequest{})
	if err != nil {
		return "", err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return resp.Hostname, nil
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
	resp, err := cli.SetHostname(context.Background(), &pb.SetHostnameRequest{Hostname: hostname})
	if err != nil {
		return err
	}
	logrus.Infof("[response]:\n%s", resp.String())
	return nil
}
