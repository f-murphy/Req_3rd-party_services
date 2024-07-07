package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"req3rdPartyServices/internal/modules"
)

func ExecuteTask(taskID int, task modules.Task) {
	if taskID == 0 {
		log.Println("incorrect task id")
	}

	switch task.Method {
	case "GET":
		executeGetTask(taskID, task)
	case "POST":
		executePostTask(taskID, task)
	case "PUT":
		executePutTask(taskID, task)
	case "DELETE":
		executeDeleteTask(taskID, task)
	}
}

func executeGetTask(taskID int, task modules.Task) {
	resp, err := http.Get(task.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

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
		Length:         fmt.Sprintf("%d", len(body)),
	}

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}

func executePostTask(taskID int, task modules.Task) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", task.Headers["Content-Type"])
	req.Header.Set("Authorization", task.Headers["Authorization"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
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
		Length:         fmt.Sprintf("%d", len(jsonTask)),
	}

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}

func executePutTask(taskID int, task modules.Task) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("PUT", task.Url, bytes.NewBuffer(jsonTask))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", task.Headers["Content-Type"])
	req.Header.Set("Authorization", task.Headers["Authorization"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
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
	}

	taskStatus := modules.TaskStatus{
		Id:             taskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
		Length:         fmt.Sprintf("%d", len(jsonTask)),
	}

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}

func executeDeleteTask(taskID int, task modules.Task) {
	req, err := http.NewRequest("DELETE", task.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", task.Headers["Authorization"])
	req.Header.Set("Content-Type", task.Headers["Content-Type"])
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatal(err)
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
	}

	taskStatus := modules.TaskStatus{
		Id:             taskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
	}

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}
