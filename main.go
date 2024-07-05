package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Task struct {
	Method  string
	Url     string
	Headers map[string]string
}

type TaskResponse struct {
	TaskID int
	Task   Task
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
	http.HandleFunc("/task/{taskID}", getTaskStatus)
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

	//go getTaskStatus(task)
	fmt.Fprintf(w, "new task ID: ")
	json.NewEncoder(w).Encode(taskID)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

func getTaskStatus(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query().Get("taskID")

	var task TaskResponse

	for _, t := range tasks {
		if fmt.Sprintf("%d", t.TaskID) == taskID {
			task = t
			break
		}
	}

	if task.TaskID == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	method := task.Task.Method
	resp, err := http.Get(task.Task.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	fmt.Println(status)

	taskStatus := TaskStatus{
		Id:             task.TaskID,
		Status:         status,
		HttpStatusCode: fmt.Sprintf("%d", resp.StatusCode),
		Headers:        headers,
		Length:         fmt.Sprintf("%d", len(body)),
	}

	jsonResp, err := json.Marshal(taskStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}
