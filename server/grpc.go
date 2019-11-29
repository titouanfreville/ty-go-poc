package server

import (
	"context"
	v1 "go_poc/api/v1"
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

	// register service
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1.NewToDoServiceServer(*dbConnectionInfo))

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
