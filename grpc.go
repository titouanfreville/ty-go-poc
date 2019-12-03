package main

import (
	"context"
	"go_poc/core"
	"go_poc/server"
)

var (
	// DbConnectionInfo information to connect to DB
	DbConnectionInfo = &core.DbConnection{}
	// APIServer api server configuration
	APIServer = &core.APIServerInfo{}
)

func getConf(dbSettings *core.DbConnection, serverSetting *core.APIServerInfo) {
	_, *dbSettings, *serverSetting = core.InitConfigFromFile("config.yml")
}

func main() {
	ctx := context.Background()
	getConf(DbConnectionInfo, APIServer)

	server.ServeGrpc(ctx, DbConnectionInfo, APIServer)
}
