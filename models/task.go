package models

type Task struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    map[string]string
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
