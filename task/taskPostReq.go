package task

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"req3rdPartyServices/models"
	"sync"
)

var lastTaskID int = 0
var mutex = &sync.Mutex{}
var tasks = make(map[int]models.TaskRequest)

func PostTask(c *gin.Context) {
	var task models.Task

	if err := c.BindJSON(&task); err != nil {
		return
	}

	mutex.Lock()
	taskID := lastTaskID + 1
	lastTaskID = taskID
	mutex.Unlock()

	taskResponse := models.TaskRequest{TaskID: taskID, Task: task}
	tasks[taskID] = taskResponse
	go executeTask(taskResponse)

	c.JSON(http.StatusOK, gin.H{
		"taskID": taskID,
	})
}
