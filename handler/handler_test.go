package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"req3rdPartyServices/models"
	mock_service "req3rdPartyServices/utils/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskHandler_CreateTask(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockTaskServiceInterface(controller)
	handler := NewTaskHandler(mockService)

	task := &models.Task{
		Method: "GET",
		Url:    "https://feeds.skynews.com/feeds/rss/business.xml",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: map[string]string{
			"key": "value",
		},
	}

	taskStatus := &models.TaskStatus{
		Status:         "200 OK",
		HttpStatusCode: "200",
		Length:         "13393",
	}

	mockService.EXPECT().CreateTask(task, taskStatus).Return(1, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	jsonBytes, err := json.Marshal(task)
	assert.NoError(t, err)

	c.Request, err = http.NewRequest("POST", "/task", bytes.NewBuffer(jsonBytes))
	assert.NoError(t, err)

	handler.CreateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, fmt.Sprintf(`{"task id":%d}`, 1), w.Body.String())
}

func TestTaskHandler_GetAllTasks(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockTaskServiceInterface(controller)
	handler := NewTaskHandler(mockService)

	tasks := []*models.TaskFromDB{
		{
			Id:             1,
			Method:         "GET",
			Url:            "https://feeds.skynews.com/feeds/rss/business.xml",
			Headers:        "{\"asd\":\"asd\"}",
			Body:           "{\"asd\":\"asd\"}",
			Status:         "200 OK",
			HttpStatusCode: "200",
			Length:         "13393",
		},
		{
			Id:             2,
			Method:         "GET",
			Url:            "https://google.com",
			Headers:        "{\"asd\":\"asd\"}",
			Body:           "{\"asd\":\"asd\"}",
			Status:         "400 Bad Request",
			HttpStatusCode: "400",
			Length:         "1555",
		},
	}

	mockService.EXPECT().GetAllTasks().Return(tasks, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("GET", "/tasks", nil)

	handler.GetAllTasks(c)
	var expectedTasks []map[string]interface{}

	for _, task := range tasks {
		expectedTask := map[string]interface{}{
			"Id":             task.Id,
			"Method":         task.Method,
			"Url":            task.Url,
			"Headers":        task.Headers,
			"Body":           task.Body,
			"Status":         task.Status,
			"HttpStatusCode": task.HttpStatusCode,
			"Length":         task.Length,
		}
		expectedTasks = append(expectedTasks, expectedTask)
	}

	jsonExpectedTasks, err := json.Marshal(expectedTasks)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(jsonExpectedTasks), w.Body.String())
}

func TestTaskHandler_GetTaskById(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	mockService := mock_service.NewMockTaskServiceInterface(controller)
	handler := NewTaskHandler(mockService)

	task := &models.TaskFromDB{
		Id:             1,
		Method:         "GET",
		Url:            "https://feeds.skynews.com/feeds/rss/business.xml",
		Headers:        "{\"asd\":\"asd\"}",
		Body:           "{\"asd\":\"asd\"}",
		Status:         "200 OK",
		HttpStatusCode: "200",
		Length:         "13393",
	}

	mockService.EXPECT().GetTaskById(1).Return(task, nil)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	c.Request, _ = http.NewRequest("GET", "/task/:id", nil)

	handler.GetTask(c)

	expectedTask := map[string]interface{}{
		"Id":             task.Id,
		"Method":         task.Method,
		"Url":            task.Url,
		"Headers":        task.Headers,
		"Body":           task.Body,
		"Status":         task.Status,
		"HttpStatusCode": task.HttpStatusCode,
		"Length":         task.Length,
	}

	jsonExpectedTask, err := json.Marshal(expectedTask)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, string(jsonExpectedTask), w.Body.String())
}
