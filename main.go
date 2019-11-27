package main

import (
	"go_poc/api"
	"go_poc/core"
)

var (
	// DbConnectionInfo information to conect to DB
	DbConnectionInfo = &core.DbConnection{}
	// APIServer api server configuration
	APIServer = &core.APIServerInfo{}
)

func getConf(dbSettings *core.DbConnection, serverSetting *core.APIServerInfo) {
	*dbSettings, *serverSetting = core.InitConfig()
}

func initAPI() error {
	if err := api.StartAPI(APIServer.Hostname, APIServer.Port, DbConnectionInfo, APIServer, false); err != nil {
		return err
	}
	return nil
}

func main() {
	getConf(DbConnectionInfo, APIServer)
	initAPI()
}
