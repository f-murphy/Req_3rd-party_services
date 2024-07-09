package task

import (
	"net/http"
)

func RouteRedirection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetTaskStatus(w, r)
	case http.MethodPost:
		PostTask(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}
