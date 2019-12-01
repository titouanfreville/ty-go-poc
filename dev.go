package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"go_poc/api"
	"go_poc/core"
	"go_poc/server"
)

var (
	// DbConnectionInfo information to conect to DB
	DbConnectionInfo = &core.DbConnection{}
	// APIServer api server configuration
	APIServer = &core.APIServerInfo{}
	log       = logrus.New()
)

func getConf(dbSettings *core.DbConnection, serverSetting *core.APIServerInfo) {
	*dbSettings, *serverSetting = core.InitConfig()
}

func main() {
	ctx := context.Background()
	getConf(DbConnectionInfo, APIServer)

	go func() {
		log.Info("In async run REST")
		err := api.StartAPI(APIServer, false)
		log.Error(err)
	}()
	server.ServeGrpc(ctx, DbConnectionInfo, APIServer)
}
