package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"req3rdPartyServices/internal/modules"
	"testing"
)

func TestTaskPostReq(t *testing.T) {
	// Arrange
	task := modules.Task{
		Method: "POST",
		Url:    "https://api.vatcomply.com/rates?base=GBP",
		Headers: map[string]string{
			"Authorization": "Basic ...",
			"Content-Type":  "application/json",
		},
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.vatcomply.com/rates", bytes.NewBuffer(jsonTask))
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
