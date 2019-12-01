package main

import (
	"go_poc/api"
	"go_poc/core"
)

var (
	// DbConnectionInfo information to connect to DB
	DbConnectionInfo = &core.DbConnection{}
	// APIServer api server configuration
	APIServer = &core.APIServerInfo{}
)

func getConf(dbSettings *core.DbConnection, serverSetting *core.APIServerInfo) {
	*dbSettings, *serverSetting = core.InitConfig()
}

func main() {
	getConf(DbConnectionInfo, APIServer)

	api.StartAPI(APIServer, false)
}
