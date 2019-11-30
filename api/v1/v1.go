package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_poc/api/v1/tasks"
	"go_poc/api/v1/users"
	"google.golang.org/grpc"
)

var (
	log = logrus.New()
)

// InitEndpoints v1
func InitEndpoints(router *gin.RouterGroup, grpcConn *grpc.ClientConn) {
	tasks.NewToDoGW(router, grpcConn)
	users.NewUserGW(router, grpcConn)
}
