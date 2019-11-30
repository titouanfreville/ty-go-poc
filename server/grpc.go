package server

import (
	"context"
	"go_poc/core"
	tasks2 "go_poc/module/tasks"
	"go_poc/module/user"
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
	tasks2.RegisterToDoServiceServer(server, tasks2.NewToDoServiceServer(*dbConnectionInfo))
	user.RegisterUserServiceServer(server, user.NewUserServiceServer(*dbConnectionInfo))

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
