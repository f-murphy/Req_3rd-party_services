package models

type Task struct {
	Method      string
	Url         string
	Headers     map[string]string
	Body        map[string]string
	HeadersJSON string
	BodyJSON    string
}

type TaskStatus struct {
	Status         string
	HttpStatusCode string
	Headers        string
	Length         string
}

type TaskFromDB struct {
	Id             int
	Method         string
	Url            string
	Headers        string
	Body           string
	Status         string
	HttpStatusCode string
	Length         string
}
