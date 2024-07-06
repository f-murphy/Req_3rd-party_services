package service

import (
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

	resp, err := http.Get(task.Url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

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
		status = "in_process"
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
