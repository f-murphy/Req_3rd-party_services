package handler

import (
	"net/http"
	"req3rdPartyServices/models"
	"req3rdPartyServices/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	service service.TaskServiceInterface
}

func NewTaskHandler(services service.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{service: services}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	err := c.BindJSON(&task)
	if err != nil {
		logrus.WithError(err).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := h.service.CreateTask(&task)
	if err != nil {
		logrus.WithError(err).Error("error creating task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task id": id})
}

func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		logrus.WithError(err).Error("error getting all tasks")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	logrus.Info("All tasks successfully received")
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
	task, err := h.service.GetTaskById(taskID)
	if err != nil {
		logrus.WithError(err).Error("error getting task")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	logrus.Info("Task by ID successfully received")
	c.JSON(http.StatusOK, task)
}
