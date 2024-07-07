package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"req3rdPartyServices/internal/modules"
	"strings"
)

func redirectionTask(taskID int, task modules.Task, w http.ResponseWriter, r *http.Request) {
	if taskID == 0 {
		log.Println("incorrect task id")
	}

	switch strings.ToUpper(task.Method) {
	case "GET":
		executeGetTask(taskID, task, w, r)
	case "POST":
		executePostTask(taskID, task, w, r)
	case "PUT":
		executePutTask(taskID, task, w, r)
	case "DELETE":
		executeDeleteTask(taskID, task, w, r)
	}
}

func executeGetTask(taskID int, task modules.Task, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(task.Url)
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

	executeTask(taskID, task, "GET", body, w, r)
}

func executePostTask(taskID int, task modules.Task, w http.ResponseWriter, r *http.Request) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	executeTask(taskID, task, "POST", jsonTask, w, r)
}

func executePutTask(taskID int, task modules.Task, w http.ResponseWriter, r *http.Request) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	executeTask(taskID, task, "PUT", jsonTask, w, r)
}

func executeDeleteTask(taskID int, task modules.Task, w http.ResponseWriter, r *http.Request) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}
	executeTask(taskID, task, "DELETE", jsonTask, w, r)
}

func executeTask(taskID int, task modules.Task, method string, body []byte, w http.ResponseWriter, r *http.Request) {

	req, err := http.NewRequest(method, task.Url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", task.Headers["Content-Type"])
	req.Header.Set("Authorization", task.Headers["Authorization"])

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

	taskStatus := modules.TaskStatus{
		Id:             taskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
	}

	if body != nil {
		taskStatus.Length = fmt.Sprintf("%d", len(body))
	}

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}
