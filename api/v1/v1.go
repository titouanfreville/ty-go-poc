package v1

import (
	"github.com/gin-gonic/gin"
	"go_poc/api/v1/tasks"
	"go_poc/api/v1/users"
)

func InitEndpoints(router *gin.RouterGroup) {
	tasks.NewGateWay(router)
	users.NewGateWay(router)
}
