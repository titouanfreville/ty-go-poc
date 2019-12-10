package api

import (
	revision "github.com/appleboy/gin-revision-middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	v1 "go_poc/api/v1"
	"net/http"
	"os"
)

var (
	log = logrus.New()
)

// newRouter initialise api server.
func newRouter() *gin.Engine {
	return gin.Default()
}

// initMiddleware initialise middleware for router
func initMiddleware(router *gin.Engine) {
	router.Use(revision.Middleware()) // inject REVISION file content to headers X-Revision
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
func initRoute(router *gin.Engine) {
	initMiddleware(router)
	basicRoutes(router)
	v1gr := router.Group("/v1")
	v1.InitEndpoints(v1gr)
}

/*
// StartAPI initialise the api with provided host and port.
func StartAPI(providedAPIServerInfo *core.APIServerInfo, test bool) error {
	// Init grpc connection
	log.Info(fmt.Sprintf("%s:%d", providedAPIServerInfo.Hostname, providedAPIServerInfo.RPCPort))
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", providedAPIServerInfo.Hostname, providedAPIServerInfo.RPCPort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	router := newRouter()
	log.Out = os.Stderr
	log.Level = logrus.DebugLevel

	/*initAuth()
	admin := router.Group("/admin")
	admin.Use(dbRequiredMiddleware())
	admin.Use(Verifier(tokenAuth))
	admin.Use(Authenticator())
	initRoute(router, conn)

	if test {
		return nil
	}

	return router.Run(fmt.Sprintf("%s:%d", providedAPIServerInfo.Hostname, providedAPIServerInfo.RESTPort))
}
*/

// GetAPIRouter initialise the api with provided host and port.
func GetAPIRouter() http.Handler {
	router := newRouter()
	log.Out = os.Stderr
	log.Level = logrus.DebugLevel

	/*initAuth()
	admin := router.Group("/admin")
	admin.Use(dbRequiredMiddleware())
	admin.Use(Verifier(tokenAuth))
	admin.Use(Authenticator())*/
	initRoute(router)

	return router
}
