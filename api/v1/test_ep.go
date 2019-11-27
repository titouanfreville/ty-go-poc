package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go_poc/core"
	"go_poc/module/test"
	"net/http"
)

var (
	log     = logrus.New()
	dbStore *core.DBStore
)

// Init v1 test endpoints
func InitTestEndPoint(router *gin.RouterGroup, apiDbStore *core.DBStore) {
	testEP := router.Group("/test")
	testEP.GET("", getAllTest)
	testEP.POST("", newTest)
	dbStore = apiDbStore
}

func getAllTest(c *gin.Context) {
	res := test.GetAll(dbStore.Db)
	c.JSON(http.StatusOK, res)
}

func newTest(c *gin.Context) {
	var testData test.Test
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&testData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&testData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	
	if err := test.Insert(&testData, dbStore.Db); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, testData)
}
