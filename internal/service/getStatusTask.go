package service

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

	for _, t := range tasks {
		if t.TaskID == taskID {
			json.NewEncoder(w).Encode(t.TaskStatus)
			return
		}
	}

	http.Error(w, "Invalid task ID", http.StatusBadRequest)
}
