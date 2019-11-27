package api

import (
	revision "github.com/appleboy/gin-revision-middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_poc/api/v1"
	"go_poc/core"
	"net/http"
	"os"
)

var (
	log      = logrus.New()
	error503 = core.NewAPIError("Database", "connection", "Database is currently in maintenance state. We are doing our best to get it back online as soon as possible.", http.StatusServiceUnavailable, "ANY")
	DbStore  = &core.DBStore{}
)

// newRouter initialise api server.
func newRouter() *gin.Engine {
	return gin.Default()
}

// initMiddleware initialise middleware for router
func initMiddleware(router *gin.Engine) {
	router.Use(revision.Middleware()) // inject REVISION file content to headers X-Revision
	router.Use(dbRequiredMiddleware())
}

// dbRequireMiddleware check if db is available
func dbRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !DbStore.AssertConnectionOrReconnect() {
			c.JSON(http.StatusServiceUnavailable, error503)
			c.Abort()
			return
		}
		c.Next()
	}
}

// recoveryHandler middleware to cleanly recover from panic
func recoveryHandler(c *gin.Context, err interface{}) {
	c.HTML(500, "error.tmpl", gin.H{
		"title": "Error",
		"err":   err,
	})
}

// basicRoutes set basic routes for the API
func basicRoutes(router *gin.Engine) {
	// router.LoadHTMLGlob("templates/*")
	router.Use(dbRequiredMiddleware())
	// swagger:route GET / Test hello
	//
	// Hello World
	//
	// 	Responses:
	//    200: generalOk
	// 	  default: genericError
	router.GET("/", func(c *gin.Context) {
		log.Info("Test")
		c.String(http.StatusOK, "Hello this is POC GoLang Tankyou API")
	})
	// swagger:route GET /ping Test ping
	//
	// Pong
	//
	// Test api ping
	//
	// 	Responses:
	//    200: generalOk
	// 	  default: genericError
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// swagger:route GET /panic Test panic
	//
	// Should result in 500
	//
	// Test panic cautching
	//
	// 	Responses:
	//    500: genericError
	// 	  default: genericError
	router.GET("/panic", func(c *gin.Context) {
		panic("C'est la panique, panique, panique. Sur le périphérique")
	})
}

// initRoute initialize all routes on correct routers
func initRoute(router *gin.Engine) { // , admin *gin.RouterGroup
	initMiddleware(router)
	basicRoutes(router)
	v1gr := router.Group("/v1")
	v1.InitTestEndPoint(v1gr, DbStore)
	// planningRoutes(router, admin *gin.RouterGroup)
}

// StartAPI initialise the api with provided host and port.
func StartAPI(hostname string, port string, DbConnectionInfo *core.DbConnection, providedAPIServerInfo *core.APIServerInfo, test bool) error {
	router := newRouter()
	log.Out = os.Stderr
	log.Level = logrus.DebugLevel

	// Init DB connection
	user := DbConnectionInfo.User
	dataBase := DbConnectionInfo.Database
	pass := DbConnectionInfo.Password
	host := DbConnectionInfo.Host
	dbPort := DbConnectionInfo.Port
	DbStore.InitConnection(user, dataBase, pass, host, dbPort)

	/*initAuth()
	admin := router.Group("/admin")
	admin.Use(dbRequiredMiddleware())
	admin.Use(Verifier(tokenAuth))
	admin.Use(Authenticator())*/
	initRoute(router)

	if test {
		return nil
	}

	return router.Run(providedAPIServerInfo.Hostname + ":" + providedAPIServerInfo.Port)
}
