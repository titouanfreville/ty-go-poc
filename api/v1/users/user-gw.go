package users

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
)

var (
	userClient UserServiceClient
	log        = logrus.New()
)

func NewUserGW(router *gin.RouterGroup, grpcConn *grpc.ClientConn) {
	userClient = NewUserServiceClient(grpcConn)

	userGW := router.Group("/users")
	userGW.GET("", getAll)
	userGW.POST("", newUser)
	// Just to make thinks readable
	/*{
		userGWId := userGW.Group("/:id")
		userGWId.GET("", getSingle)
		userGWId.PUT("", update)
		userGWId.PATCH("", update)
	}*/
}

func getAll(c *gin.Context) {
	res, err := userClient.ReadAll(context.Background(), &ReadAllRequest{Api: "v1"})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func newUser(c *gin.Context) {
	var userData User
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&userData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&userData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	res, err := userClient.Create(context.Background(), &CreateRequest{Api: apiVersion, User: &userData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}
