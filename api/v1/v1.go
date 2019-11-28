package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_poc/core"
)

var (
	log     = logrus.New()
	dbStore *core.DBStore
)

// InitEndpoints v1
func InitEndpoints(router *gin.RouterGroup, apiDbStore *core.DBStore) {
	dbStore = apiDbStore
	initUserEndPoint(router)
	initTaskEndPoint(router)
}
