package service

import (
	"req3rdPartyServices/models"
	"req3rdPartyServices/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TaskServiceInterface interface {
	CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error)
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskService struct {
	repo repository.TaskRepositoryInterface
}

func NewTaskService(repo repository.TaskRepositoryInterface) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error) {
	return s.repo.CreateTask(task, taskStatus)
}

func (s *TaskService) GetAllTasks() ([]*models.TaskFromDB, error) {
	return s.repo.GetAllTasks()
}

func (s *TaskService) GetTaskById(id int) (*models.TaskFromDB, error) {
	return s.repo.GetTaskById(id)
}
