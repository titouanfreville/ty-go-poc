package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log = logrus.New()
)

// InitEndpoints v1
func InitEndpoints(router *gin.RouterGroup, grpcConn *grpc.ClientConn) {
	NewToDoGW(router, grpcConn)
	NewUserGW(router, grpcConn)
}
