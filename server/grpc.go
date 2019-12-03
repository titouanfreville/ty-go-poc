package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go_poc/api/v1/tasks"
	"go_poc/api/v1/users"
	"go_poc/core"
	"net"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"
)

var (
	log = logrus.New()
)

func ApiVersionMiddleWareU(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	log.Info(start)
	log.Info(req)
	log.Info(*info)
	return handler(ctx, req)
}

// RunServer runs gRPC service to publish ToDo service
func ServeGrpc(ctx context.Context, dbConnectionInfo *core.DbConnection, apiServer *core.APIServerInfo) error {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", apiServer.Hostname, apiServer.RPCPort))
	if err != nil {
		return err
	}

	// register services
	server := grpc.NewServer(grpc.UnaryInterceptor(ApiVersionMiddleWareU))
	tasks.RegisterToDoServiceServer(server, tasks.NewToDoServiceServer(*dbConnectionInfo))
	users.RegisterUserServiceServer(server, users.NewUserServiceServer(*dbConnectionInfo))

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Printf("starting gRPC server on  %s:%d", apiServer.Hostname, apiServer.RPCPort)
	return server.Serve(listen)
}
