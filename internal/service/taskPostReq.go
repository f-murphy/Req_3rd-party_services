package service

import (
	"encoding/json"
	"net/http"
	"req3rdPartyServices/internal/modules"
	"strings"
)

var lastTaskID int = 0

var tasks []modules.TaskResponse

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task modules.Task

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

	taskID := lastTaskID + 1
	lastTaskID = taskID

	taskResponse := modules.TaskResponse{TaskID: taskID, Task: task}
	tasks = append(tasks, taskResponse)

	go redirectionTask(taskID, task)

	json.NewEncoder(w).Encode(taskID)
}
