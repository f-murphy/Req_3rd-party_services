package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Task struct {
	Method  string
	Url     string
	Headers map[string]string
}

type TaskResponse struct {
	TaskID     int
	Task       Task
	TaskStatus any
}

type TaskStatus struct {
	Id             int
	Status         string
	HttpStatusCode string
	Headers        map[string]string
	Length         string
}

var lastTaskID int = 0

var tasks []TaskResponse

func main() {
	http.HandleFunc("/tasks", taskHandler)
	startServer()
}

func startServer() {
	fmt.Println("Listening on port 8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskStatus(w, r)
	case http.MethodPost:
		postTask(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	taskID := lastTaskID + 1
	lastTaskID = taskID

	taskResponse := TaskResponse{TaskID: taskID, Task: task}
	tasks = append(tasks, taskResponse)

	go executeTask(taskID, task)

	fmt.Fprintf(w, "new task ID: ")
	json.NewEncoder(w).Encode(taskID)
}

func executeTask(taskID int, task Task) {
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

	taskStatus := TaskStatus{
		Id:             taskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
		Length:         fmt.Sprintf("%d", len(body)),
	}
	fmt.Println(status)

	for i, t := range tasks {
		if t.TaskID == taskID {
			tasks[i].TaskStatus = taskStatus
			break
		}
	}
}

func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Query().Get("taskID")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, t := range tasks {
		if t.TaskID == taskID {
			json.NewEncoder(w).Encode(t.TaskStatus)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}
