package tasks

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
)

var (
	todoClient ToDoServiceClient
	log        = logrus.New()
)

func NewToDoGW(router *gin.RouterGroup, grpcConn *grpc.ClientConn) {
	todoClient = NewToDoServiceClient(grpcConn)

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
	res, err := todoClient.ReadAll(context.Background(), &ReadAllRequest{Api: "v1"})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
	res, err := todoClient.Read(context.Background(), &ReadRequest{Api: "v1", Id: id})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func newTask(c *gin.Context) {
	var taskData ToDo
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&taskData); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	res, err := todoClient.Create(context.Background(), &CreateRequest{Api: "v1", ToDo: &taskData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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

	res, err := todoClient.Read(context.Background(), &ReadRequest{Api: "v1", Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, "Cat not found")
		return
	}
	taskData := res.ToDo
	if c.ContentType() == "multipart/form-data" {
		if err := c.Bind(&taskData); err != nil {
			log.Info(err)
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err := c.BindJSON(&taskData); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	_, err = todoClient.Update(context.Background(), &UpdateRequest{Api: "v1", ToDo: taskData})
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, taskData)
}
