package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"req3rdPartyServices/models"
	"req3rdPartyServices/service"
	"strconv"
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

	taskStatusChan := make(chan *models.TaskStatus)
	errChan := make(chan error)

	go func() {
		taskStatus, err := service.ExecuteTask(&task)
		if err != nil {
			errChan <- err
			return
		}
		taskStatusChan <- taskStatus
	}()

	select {
	case taskStatus := <-taskStatusChan:
		id, err := h.service.CreateTask(&task, taskStatus)
		if err != nil {
			logrus.WithError(err).Error("error creating task in DB")
			c.JSON(http.StatusInternalServerError, gin.H{"error creating task in DB": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task id": id})
	case err := <-errChan:
		logrus.WithError(err).Error("error executing task")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
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
	task, err := h.service.GetTaskById(taskID)
	if err != nil {
		logrus.WithError(err).Error("error getting task")
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
