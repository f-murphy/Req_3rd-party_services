package handler

import (
	"net/http"
	"req3rdPartyServices/models"
	"req3rdPartyServices/service"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var mutex = &sync.Mutex{}

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
	logrus.Info("Task binding successfully")

	taskStatusChan := make(chan *models.TaskStatus)
	errChan := make(chan error)

	go func(task *models.Task, taskStatusChan chan *models.TaskStatus, errChan chan error) {
		defer func() {
			mutex.Lock()
            close(taskStatusChan)
            close(errChan)
            mutex.Unlock()
		}()

		taskStatus, err := service.ExecuteTask(task)
		if err != nil {
			mutex.Lock()
			errChan <- err
			mutex.Unlock()
			return
		}
		mutex.Lock()
		taskStatusChan <- taskStatus
		mutex.Unlock()
	}(&task, taskStatusChan, errChan)

	select {
	case taskStatus, ok := <-taskStatusChan:
		if !ok {
			return
		}
		id, err := h.service.CreateTask(&task, taskStatus)
		if err != nil {
			logrus.WithError(err).Error("error creating task in DB")
			c.JSON(http.StatusInternalServerError, gin.H{"error creating task in DB": err.Error()})
			return
		}
		logrus.Info("The task has been successfully created in the database")
		c.JSON(http.StatusOK, gin.H{"task id": id})
	case err, ok := <-errChan:
		if !ok {
			return
		}
		logrus.WithError(err).Error("error executing task")
		c.JSON(http.StatusBadRequest, gin.H{"error:": err})
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
	logrus.Info("All tasks successfully received")
	c.JSON(http.StatusOK, tasks)
}

// @Summary Get task
// @Description Get task by id
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.TaskFromDB
// @Failure 400
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
	logrus.Info("Task by ID successfully received")
	c.JSON(http.StatusOK, task)
}
