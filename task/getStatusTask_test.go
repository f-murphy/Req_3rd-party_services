package task

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTaskStatus(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "?taskID=asdasd", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	// Act
	GetTaskStatus(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code 400, but got %d", w.Code)
	}
	if w.Body.String() != "Invalid task ID\n" {
		t.Errorf("expected response body to be 'Invalid task ID', but got '%s'", w.Body.String())
	}
}
