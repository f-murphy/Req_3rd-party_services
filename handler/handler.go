package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"req3rdPartyServices/executor"
	"req3rdPartyServices/models"
	"req3rdPartyServices/service"
	"strconv"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(services *service.TaskService) *TaskHandler {
	return &TaskHandler{service: services}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task

	err := c.BindJSON(&task)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskStatus, err := executor.ExecuteTask(task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err,
		})
		return
	}

	err = h.service.CreateTask(&task, taskStatus)
	if err != nil {
		c.JSON(500, gin.H{"error while send in DB": err.Error()})
		return
	}

	c.JSON(201, gin.H{"task": task, "taskStatus": taskStatus})
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}
	task, err := h.service.GetTask(taskID)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, task)
}
