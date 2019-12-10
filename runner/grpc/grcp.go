package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go_poc/api/v1/tasks"
	"go_poc/api/v1/users"
	"go_poc/core"

	micro "github.com/micro/go-micro"
)

var (
	dbConf core.DbConnection
	log    = logrus.New()
)

func getConfig() {
	_, dbConf, _ = core.InitConfigFromFile("config.yml")
}

func main() {
	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("grpc"),
		micro.WrapHandler(core.EnforceVersion),
		micro.Version("v1"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register handler
	if err := tasks.RegisterToDoServiceHandler(service.Server(), tasks.NewHandler(dbConf)); err != nil {
		log.Error(err)
	}
	if err := users.RegisterUserServiceHandler(service.Server(), users.NewHandler(dbConf)); err != nil {
		log.Error(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
