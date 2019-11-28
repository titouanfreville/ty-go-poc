package v1

import (
	"github.com/gin-gonic/gin"
	"go_poc/module/tasks"
	"net/http"
	"strconv"
)

// Init v1 test endpoints
func initTaskEndPoint(router *gin.RouterGroup) {
	testEP := router.Group("/task")
	testEP.GET("", getAllTask)
	testEP.POST("", newTask)
	testEP.PUT("/:id", updateTask)
}

func getAllTask(c *gin.Context) {
	res := tasks.GetAll(dbStore.Db)
	c.JSON(http.StatusOK, res)
}

func newTask(c *gin.Context) {
	var taskData tasks.Task
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if err := tasks.Insert(&taskData, dbStore.Db); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, taskData)
}

func updateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ID is not an integer")
		return
	}

	taskData := tasks.GetOne(id, dbStore.Db)
	if taskData == nil {
		c.JSON(http.StatusNotFound, "Cat not found")
		return
	}

	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	if err := tasks.Update(taskData, dbStore.Db); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusCreated, taskData)
}
