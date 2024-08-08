package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"req3rdPartyServices/models"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var mutex = &sync.Mutex{}

func ExecuteTask(task *models.Task) (taskStatus *models.TaskStatus, err error) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		logrus.WithError(err).Error("error marshaling task")
		return nil, err
	}

	req, err := http.NewRequest(strings.ToUpper(task.Method), task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		logrus.WithError(err).Error("error creating new request")
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("error sending request")
		return nil, err
	}
	defer func() {
		mutex.Lock()
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Error("error closing response body")
		}
		mutex.Unlock()
	}()

	taskStatus = &models.TaskStatus{
		Status:         resp.Status,
		HttpStatusCode: strconv.Itoa(resp.StatusCode),
		Length:         strconv.FormatInt(resp.ContentLength, 10),
	}

	return taskStatus, err
}
