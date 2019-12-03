package core

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

// YMLConfig config struct for .yml config file
type YMLConfig struct {
	Database DbConnection  `yaml:"database,flow"`
	Server   APIServerInfo `yaml:"server,flow"`
	// Others ...
}

// DbConnection information to connect to DB
type DbConnection struct {
	User            string `yaml:"user"`
	Database        string `yaml:"name"`
	Password        string `yaml:"password"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DefaultTimeZone string `yaml:"tz"`
}

// APIServerInfo information on API server
type APIServerInfo struct {
	Hostname     string `yaml:"host"`
	RPCPort      int    `yaml:"rpc_port"`
	RESTPort     int    `yaml:"rest_port"`
	JWTSecretKey string `yaml:"jwt_sk"`
	LogLevel     string `yaml:"log_level"`
}

// InitConfig get configuration for project
func InitConfig() (DbConnection, APIServerInfo) {
	// Default configurations
	dbConnection := DbConnection{
		User:            "tankyou_poc",
		Database:        "tankyou_poc",
		Password:        "tankyou_poc",
		Host:            "0.0.0.0",
		Port:            5432,
		DefaultTimeZone: "Europe/Paris",
	}
	APIServer := APIServerInfo{
		Hostname:     "0.0.0.0",
		RESTPort:     3000,
		RPCPort:      3001,
		JWTSecretKey: "MagicalTokenIsTheBest",
	}

	// Default host for DB in Docker containers
	if os.Getenv("ENVTYPE") == "container" {
		log.Print("<><><><> Setting host to container default \n")
		dbConnection.Host = "database"
	}
	// Get values set in env
	if apiPort := os.Getenv("API_PORT"); apiPort != "" {
		log.Print("<><><><> Setting api port \n")
		restPort, err := strconv.Atoi(apiPort)
		if err == nil {
			APIServer.RESTPort = restPort
		}
	}
	if apiHostname := os.Getenv("API_HOST"); apiHostname != "" {
		log.Print("<><><><> Setting api hostname \n")
		APIServer.Hostname = apiHostname
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		log.Print("<><><><> Setting JWT secret \n")
		APIServer.JWTSecretKey = jwtSecret
	}
	// Will be erased if user is not root
	if dbRootPassword := os.Getenv("MYSQL_ROOT_PASSWORD"); dbRootPassword != "" {
		log.Print("<><><><> Setting db root password \n")
		dbConnection.Password = dbRootPassword
	}
	if dbUser := os.Getenv("MYSQL_USER"); dbUser != "" {
		log.Print("<><><><> Setting db user and user password \n")
		dbConnection.User = dbUser
		// Can be empty. Should be define when user is define
		dbConnection.Password = os.Getenv("MYSQL_PASSWORD")
	}
	if dbName := os.Getenv("MYSQL_DATABASE"); dbName != "" {
		log.Print("<><><><> Setting db name \n")
		dbConnection.Database = dbName
	}
	if dbPort := os.Getenv("MYSQL_PORT"); dbPort != "" {
		log.Print("<><><><> Setting db port \n")
		dataBPort, err := strconv.Atoi(dbPort)
		if err == nil {
			dbConnection.Port = dataBPort
		}
	}
	if dbHost := os.Getenv("MYSQL_HOST"); dbHost != "" {
		log.Print("<><><><> Setting db host \n")
		dbConnection.Host = dbHost
	}
	if defTimeZone := os.Getenv("DEFAULT_TIME_ZONE"); defTimeZone != "" {
		log.Print("<><><><> Setting db host \n")
		dbConnection.DefaultTimeZone = defTimeZone
	}

	// Return new configs
	return dbConnection, APIServer
}

func InitConfigFromFile(filepath string) (bool, DbConnection, APIServerInfo) {
	var config YMLConfig

	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Errorf("Couldn't read config file at %s", filepath)
		return false, DbConnection{}, APIServerInfo{}
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Errorf("Couldn't read config file at %s", filepath)
		return false, DbConnection{}, APIServerInfo{}
	}

	return true, config.Database, config.Server
}
