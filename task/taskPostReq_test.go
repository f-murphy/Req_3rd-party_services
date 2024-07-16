package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"req3rdPartyServices/models"
	"testing"
)

func TestTaskPostReq(t *testing.T) {
	// Arrange
	task := models.Task{
		Method: "POST",
		Url:    "https://api.vatcomply.com/rates?base=GBP",
		Headers: map[string]string{
			"Authorization": "Basic ...",
			"Content-Type":  "application/json",
		},
		Body: map[string]string{
			"asd": "asd",
		},
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	// Act
	PostTask(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status code 200, but got %d", w.Code)
	}

	var taskID int
	err = json.NewDecoder(w.Body).Decode(&taskID)
	if err != nil {
		t.Errorf("expected taskID in response, but got error %s", err.Error())
	}
}

func TestTaskPostReq_noBody(t *testing.T) {
	task := models.Task{
		Method: "POST",
		Url:    "https://api.vatcomply.com/rates?base=GBP",
		Headers: map[string]string{
			"Authorization": "Basic ...",
			"Content-Type":  "application/json",
		},
		Body: map[string]string{},
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	PostTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	if w.Body.String() != "No body for request\n" {
		t.Errorf("expected error message 'No body for request', got '%s'", w.Body.String())
	}
}
