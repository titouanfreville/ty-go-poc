package v1

import (
	"github.com/gin-gonic/gin"
	"go_poc/module/user"
	"net/http"
)

// Init v1 test endpoints
func initUserEndPoint(router *gin.RouterGroup) {
	testEP := router.Group("/user")
	testEP.GET("", getAllUser)
	testEP.POST("", newUser)
}

func getAllUser(c *gin.Context) {
	res := user.GetAll(dbStore.Db)
	c.JSON(http.StatusOK, res)
}

func newUser(c *gin.Context) {
	var userData user.User
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

	if err := user.Insert(&userData, dbStore.Db); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, userData)
}
