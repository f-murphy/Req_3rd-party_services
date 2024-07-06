package modules

type TaskStatus struct {
	Id             int
	Status         string
	HttpStatusCode string
	Headers        map[string]string
	Length         string
}
