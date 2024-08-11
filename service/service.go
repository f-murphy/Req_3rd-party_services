package service

import (
	"encoding/json"
	"req3rdPartyServices/models"
	"req3rdPartyServices/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TaskServiceInterface interface {
	CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error)
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskService struct {
	repo       repository.TaskRepositoryInterface
	redis      *redis.Client
	expiration time.Duration
}

func NewTaskService(repo repository.TaskRepositoryInterface, redis *redis.Client, expiration time.Duration) *TaskService {
	return &TaskService{repo: repo, redis: redis, expiration: expiration}
}

func (s *TaskService) CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error) {
	jsonHeaders, err := json.Marshal(task.Headers)
	if err != nil {
		return 0, err
	}

	jsonBody, err := json.Marshal(task.Body)
	if err != nil {
		return 0, err
	}

	task.HeadersJSON = string(jsonHeaders)
	task.BodyJSON = string(jsonBody)

	return s.repo.CreateTask(task, taskStatus)
}

func (s *TaskService) GetAllTasks() ([]*models.TaskFromDB, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTaskById(id int) (*models.TaskFromDB, error) {
	return s.repo.GetTaskById(id)
}
