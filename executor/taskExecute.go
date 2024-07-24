package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"req3rdPartyServices/models"
	"strings"
)

func ExecuteTask(task models.Task) (taskStatus *models.TaskStatus, err error) {
	jsonTask, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(strings.ToUpper(task.Method), task.Url, bytes.NewBuffer(jsonTask))
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

	taskStatus = &models.TaskStatus{
		Status:         resp.Status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Length:         fmt.Sprintf("%d", resp.ContentLength),
	}

	return taskStatus, err

}
