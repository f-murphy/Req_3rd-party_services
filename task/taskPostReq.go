package task

import (
	"encoding/json"
	"net/http"
	"req3rdPartyServices/models"
	"strings"
	"sync"
)

var lastTaskID int = 0
var mutex = &sync.Mutex{}
var tasks = make(map[int]models.TaskResponse)

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if (strings.ToUpper(task.Method) == "POST" || strings.ToUpper(task.Method) == "PUT" || strings.ToUpper(task.Method) == "DELETE") && task.Body == nil {
		http.Error(w, "No body for request", http.StatusBadRequest)
		return
	}

	method := strings.ToUpper(task.Method)
	allowedMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}

	if !allowedMethods[method] {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	taskID := lastTaskID + 1
	lastTaskID = taskID
	mutex.Unlock()

	taskResponse := models.TaskResponse{TaskID: taskID, Task: task}
	tasks[taskID] = taskResponse
	go redirectionTask(taskResponse)

	json.NewEncoder(w).Encode(taskID)
}
