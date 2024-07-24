package service

import (
	"req3rdPartyServices/models"
	"req3rdPartyServices/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task *models.Task, taskStatus *models.TaskStatus) error {
	return s.repo.CreateTask(task, taskStatus)
}

func (s *TaskService) GetTask(id int) (*models.TaskFromDB, error) {
	return s.repo.GetTask(id)
}
