package users

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

var (
	log     = logrus.New()
	service UserService
)

func NewGateWay(router *gin.RouterGroup) {
	rest := "rest"
	service = NewClient(&rest, nil)

	userGW := router.Group("/users")
	userGW.GET("", getAll)
	userGW.POST("", newUser)
	// Just to make thinks readable
	{
		userGWId := userGW.Group("/:id")
		userGWId.GET("", getSingle)
		// userGWId.PUT("", update)
		// userGWId.PATCH("", update)
	}
}

func getAll(c *gin.Context) {
	res, err := service.ReadAll(c, &ReadAllRequest{Api: "v1"})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func getSingle(c *gin.Context) {
	id, argErr := strconv.ParseInt(c.Param("id"), 10, 64)
	if argErr != nil {
		c.JSON(http.StatusBadRequest, "ID is not an integer")
		return
	}
	res, err := service.Read(c, &ReadRequest{Api: "v1", Id: id})
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
	res, err := service.Create(c, &CreateRequest{Api: "v1", User: &userData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}
