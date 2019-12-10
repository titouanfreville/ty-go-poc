package main

import (
	web "github.com/micro/go-micro/web"
	"github.com/sirupsen/logrus"
	"go_poc/api"
)

var (
	log = logrus.New()
)

func main() {
	// Create a new service. Optionally include some options here.
	service := web.NewService(
		web.Name("greeter.com"),
		web.Handler(api.GetAPIRouter()),
		web.Address("localhost:8081"),
		web.Version("v1"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
