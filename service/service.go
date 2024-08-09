package service

import (
	"context"
	"encoding/json"
	"fmt"
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
	repo    repository.TaskRepositoryInterface
	redis   *redis.Client
	cacheTTL time.Duration
}

func NewTaskService(repo repository.TaskRepositoryInterface, redis *redis.Client, cacheTTL time.Duration) *TaskService {
	return &TaskService{repo: repo, redis: redis, cacheTTL: cacheTTL}
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
	cacheKey := "tasks_all"
	ctx := context.Background()

	cache, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var tasks []* models.TaskFromDB
		err = json.Unmarshal([]byte(cache), &tasks)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}

	tasks, err := s.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}

	jsonTasks, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}

	err = s.redis.Set(ctx, cacheKey, jsonTasks, s.cacheTTL).Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskService) GetTaskById(id int) (*models.TaskFromDB, error) {
	cacheKey := fmt.Sprintf("task_%d", id)
	ctx := context.Background()

	cache, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var task *models.TaskFromDB
		err = json.Unmarshal([]byte(cache), &task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}

	task, err := s.repo.GetTaskById(id)
	if err != nil {
		return nil, err
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = s.redis.Set(ctx, cacheKey, jsonTask, s.cacheTTL).Err()
	if err != nil {
		return nil, err
	}

	return task, nil
}
