package main

import (
	"fmt"
	"net"
	"os"
	"time"

	pb "github.com/kata-containers/agent/protocols/grpc"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/yamux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)


var (
	channelExistMaxTries   = 200
	channelExistWaitTime   = 50 * time.Millisecond
	channelCloseTimeout    = 5 * time.Second

	// tracing enables opentracing support
	tracing = true

	grpcContext context.Context
)


type sandbox struct {
	channel    serialChannel
}

func main() {
	var (
		err error
		servErr error
		grpcServer *grpc.Server
	 	serverOpts []grpc.ServerOption
	)

	s := sandbox{}
	serverOpts = append(serverOpts, grpc.UnaryInterceptor(makeUnaryInterceptor()))
	grpcServer = grpc.NewServer(serverOpts...)

	pb.RegisterAgentServiceServer(grpcServer, grpcImpl)
	pb.RegisterHealthServer(grpcServer, grpcImpl)
	s.server = grpcServer


	for {
		logrus.Info("agent grpc server starts")

		err = s.channel.setup()
		if err != nil {
			logrus.WithError(err).Warn("Failed to setup agent grpc channel")
			return
		}
		logrus.Info("channel setup complete")

		err = s.channel.wait()
		if err != nil {
			logrus.WithError(err).Warn("Failed to wait agent grpc channel ready")
			return
		}
		logrus.Info("channel wait complete")

		var l net.Listener
		l, err = s.channel.listen()
		if err != nil {
			logrus.WithError(err).Warn("Failed to create agent grpc listener")
			return
		}
		logrus.Info("channel listen complete")

		// l is closed when Serve() returns
		servErr = grpcServer.Serve(l)
		if servErr != nil {
			logrus.WithError(servErr).Warn("agent grpc server quits")
		}
		logrus.Info("channel serve complete")

		err = s.channel.teardown()
		if err != nil {
			logrus.WithError(err).Warn("agent grpc channel teardown failed")
		}
		logrus.Info("channel teardown complete")

		// Based on the definition of grpc.Serve(), the function
		// returns nil in case of a proper stop triggered by either
		// grpc.GracefulStop() or grpc.Stop(). Those calls can only
		// be issued by the chain of events coming from DestroySandbox
		// and explicitly means the server should not try to listen
		// again, as the sandbox is being completely removed.
		if servErr == nil {
			logrus.Info("agent grpc server has been explicitly stopped")
			return
		}
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

func (c *serialChannel) wait() error {
	return nil
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

func getGRPCContext() context.Context {
	if grpcContext != nil {
		return grpcContext
	}

	logrus.Warn("Creating gRPC context as none found")

	return context.Background()
}

func makeUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(origCtx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
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

		if !tracing {
			// Just log call details
			elapsed = time.Since(start)
			message = resp.(proto.Message)

			logger := logrus.WithFields(logrus.Fields{
				"request":  info.FullMethod,
				"duration": elapsed.String(),
				"resp":     message.String()})
			logger.Debug("request end")
		}

		return resp, err
	}
}