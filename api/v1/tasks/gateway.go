package tasks

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/util/log"
	"net/http"
	"strconv"
)

var (
	service ToDoService
)

// NewGateWay initialize GateWay for task services
func NewGateWay(router *gin.RouterGroup) {
	rest := "rest"
	service = NewClient(&rest, nil)
	todoGW := router.Group("/tasks")
	todoGW.GET("", getAllTask)
	todoGW.POST("", newTask)
	// Just to make thinks readable
	{
		todoGWId := todoGW.Group("/:id")
		todoGWId.GET("", getSingleTask)
		todoGWId.PUT("", updateTask)
		todoGWId.PATCH("", updateTask)
	}
}

func getAllTask(c *gin.Context) {
	c.Set("service_version", "v1")
	res, err := service.ReadAll(c, &ReadAllRequest{Api: "v1"})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func getSingleTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ID is not an integer")
		return
	}
	res, err := service.Read(c, &ReadRequest{Api: "v1", Id: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func newTask(c *gin.Context) {
	var taskData ToDo
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := c.BindJSON(&taskData); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	res, err := service.Create(c, &CreateRequest{Api: "v1", ToDo: &taskData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func updateTask(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "ID is not an integer")
		return
	}

	res, err := service.Read(c, &ReadRequest{Api: "v1", Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, "Cat not found")
		return
	}
	taskData := res.ToDo
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&taskData); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := c.BindJSON(&taskData); err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
	}

	_, err = service.Update(c, &UpdateRequest{Api: "v1", ToDo: taskData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, taskData)
}
