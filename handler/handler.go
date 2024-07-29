package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
		logrus.WithError(err).Error("error binding JSON")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	taskStatus, err := executor.ExecuteTask(task)
	if err != nil {
		logrus.WithError(err).Error("error executing task")
		c.JSON(http.StatusBadRequest, gin.H{
			"error:": err,
		})
		return
	}

	err = h.service.CreateTask(&task, taskStatus)
	if err != nil {
		logrus.WithError(err).Error("error creating task in DB")
		c.JSON(500, gin.H{"error creating task in DB": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": task, "taskStatus": taskStatus})
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		logrus.WithError(err).Error("error getting all tasks")
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	taskID, err := strconv.Atoi(id)

	if err != nil {
		logrus.WithError(err).Error("error parsing task ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
		})
		return
	}
	task, err := h.service.GetTask(taskID)
	if err != nil {
		logrus.WithError(err).Error("error getting task")
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
