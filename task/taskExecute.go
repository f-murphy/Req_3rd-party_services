package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"req3rdPartyServices/models"
	"strings"
)

func redirectionTask(task models.TaskResponse) {
	switch strings.ToUpper(task.Task.Method) {
	case "GET":
		executeGetTask(task)
	case "POST":
		executePostTask(task)
	case "PUT":
		executePutTask(task)
	case "DELETE":
		executeDeleteTask(task)
	}
}

func executeGetTask(task models.TaskResponse) {
	resp, err := http.Get(task.Task.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	executeTask(task, "GET", body)
}

func executePostTask(task models.TaskResponse) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	executeTask(task, "POST", jsonTask)
}

func executePutTask(task models.TaskResponse) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	executeTask(task, "PUT", jsonTask)
}

func executeDeleteTask(task models.TaskResponse) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	executeTask(task, "DELETE", jsonTask)
}

func executeTask(task models.TaskResponse, method string, body []byte) {
	req, err := http.NewRequest(method, task.Task.Url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	headers := make(map[string]string)
	for k, v := range resp.Header {
		headers[k] = v[0]
	}

	var status string
	switch resp.StatusCode {
	case http.StatusOK:
		status = "done"
	case http.StatusAccepted:
		status = "in process"
	case http.StatusMethodNotAllowed:
		status = "incorrect method"
	case http.StatusInternalServerError:
		status = "error"
	default:
		status = "new"
	}

	taskStatus := models.TaskStatus{
		Id:             task.TaskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
	}

	if body != nil {
		taskStatus.Length = fmt.Sprintf("%d", len(body))
	}

	taskResp, ok := tasks[task.TaskID]
	if !ok {
		log.Println("Task not found in tasks map")
		return
	}

	taskResp.TaskStatus = taskStatus

	tasks[task.TaskID] = taskResp
}
