package server

import (
	"context"
	"go_poc/api/v1/tasks"
	"go_poc/api/v1/users"
	"go_poc/core"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish ToDo service
func ServeGrpc(ctx context.Context, dbConnectionInfo *core.DbConnection, apiServer *core.APIServerInfo) error {
	listen, err := net.Listen("tcp", apiServer.Hostname+":"+apiServer.RPCPort)
	if err != nil {
		return err
	}

	// register services
	server := grpc.NewServer()
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
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
