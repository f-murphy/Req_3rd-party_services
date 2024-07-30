package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"req3rdPartyServices/configs"
	"req3rdPartyServices/models"
	"req3rdPartyServices/repository"
	"req3rdPartyServices/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var taskService *service.TaskService
var taskHandler *TaskHandler

func init() {

	if err := configs.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing configs")
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize db")
	}

	repos := repository.NewTaskRepository(db)
	taskService = service.NewTaskService(repos)
	taskHandler = NewTaskHandler(taskService)
}

func Test_CreateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	tasks := []models.Task{
		{
			Method:  "GET",
			Url:     "http://google.com",
			Headers: map[string]string{},
			Body:    map[string]string{},
		},
		{
			Method:  "POST",
			Url:     "https://feeds.skynews.com/feeds/rss/business.xml",
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    map[string]string{},
		},
		{
			Method:  "GET",
			Url:     "https://fesadfsadfeds.skynews.com/feeds/rss/business.xml",
			Headers: map[string]string{},
			Body:    map[string]string{"id": "2"},
		},
	}

	for _, test := range tasks {
		jsonTask, err := json.Marshal(test)
		if err != nil {
			t.Errorf("error during marshal task")
		}

		c.Request = httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
		c.Request.Header.Set("Content-Type", "application/json")

		taskHandler.CreateTask(c)

		if w.Code != http.StatusOK {
			t.Errorf("Bad status %d during test CreateTask", w.Code)
		}
	}
}

func Test_GetAllTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/tasks", nil)

	taskHandler.GetAllTasks(c)

	if w.Code != http.StatusOK {
		t.Errorf("Bad status %d during test GetAllTasks", w.Code)
	}
}

func Test_GetTaskById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/tasks/1", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	taskHandler.GetTask(c)

	if w.Code != http.StatusOK {
		t.Errorf("Bad status %d during test GetTask", w.Code)
	}
}

func Test_GetTaskById_ParseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "/tasks/abc", nil)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}

	taskHandler.GetTask(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Bad status %d during test GetTask with parse error", w.Code)
	}
}
