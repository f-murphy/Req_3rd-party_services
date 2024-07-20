package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"req3rdPartyServices/models"
	"strings"
)

var headers = make(map[string]string)

func executeTask(task models.TaskRequest) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(strings.ToUpper(task.Task.Method), task.Task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	taskStatus := models.TaskStatus {
		Id: task.TaskID,
		Status: resp.Status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers: headers,
		Length: fmt.Sprintf("%d", resp.ContentLength),
	}

	taskResp, ok := tasks[task.TaskID]
	if !ok {
		log.Println("Task not found in tasks map")
	}

	taskResp.TaskStatus = taskStatus
	tasks[task.TaskID] = taskResp

}
