package models

type Task struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    map[string]string
}