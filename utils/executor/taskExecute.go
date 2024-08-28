package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"req3rdPartyServices/metrics"
	"req3rdPartyServices/models"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func ExecuteTask(task *models.Task) (taskStatus *models.TaskStatus, err error) {
	startTime := time.Now()
	defer func() {
		metrics.TaskExecutionDuration.Observe(time.Since(startTime).Seconds())
	}()

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
		if err := resp.Body.Close(); err != nil {
			logrus.WithError(err).Error("error closing response body")
		}
	}()

	taskStatus = &models.TaskStatus{
		Status:         resp.Status,
		HttpStatusCode: strconv.Itoa(resp.StatusCode),
		Length:         strconv.FormatInt(resp.ContentLength, 10),
	}

	if err != nil {
		metrics.TaskExecutionErrorsTotal.WithLabelValues(err.Error()).Inc()
	}
	return taskStatus, err
}
