package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"go_poc/api/v1/greeter"
	"time"
)

func main() {
	// Create a new service
	service := micro.NewService(micro.Name("greeter.client"))
	// Initialise the client and parse command line flags
	service.Init()
	// Create new greeter client
	grServ := greeter.NewGreeterService("greeter", service.Client())

	// Call the greeter
	rsp, err := grServ.Hello(context.TODO(), &greeter.Request{Name: "John"})
	time.Sleep(1)
	rsp, err = grServ.Hello(context.TODO(), &greeter.Request{Name: "John"})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print response
	fmt.Println(rsp)
}
