package task

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if val, ok := tasks[taskID]; ok {
		json.NewEncoder(w).Encode(val.TaskStatus)
	} else {
		http.Error(w, "non-existent task", http.StatusBadRequest)
	}
}
