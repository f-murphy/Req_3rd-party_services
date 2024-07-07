package service

import (
	"encoding/json"
	"net/http"
	"req3rdPartyServices/internal/modules"
)

var lastTaskID int = 0

var tasks []modules.TaskResponse

func PostTask(w http.ResponseWriter, r *http.Request) {
	var task modules.Task
	err := json.NewDecoder(r.Body).Decode(&task)
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
