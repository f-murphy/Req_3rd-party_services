package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"req3rdPartyServices/models"
	"req3rdPartyServices/repository"
	utils "req3rdPartyServices/utils/executor"
	workerpool "req3rdPartyServices/utils/worker"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type TaskServiceInterface interface {
	CreateTask(task *models.Task) (int, error)
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskService struct {
	repo       repository.TaskRepositoryInterface
	redis      *redis.Client
	expiration time.Duration
	wp         *workerpool.WorkerPool
}

func NewTaskService(
	repo repository.TaskRepositoryInterface,
	redis *redis.Client,
	expiration time.Duration,
) *TaskService {
	workerpool_size := workerpool.NewWorkerPool(10)
	s := &TaskService{
		repo:       repo,
		redis:      redis,
		expiration: expiration,
		wp:         workerpool_size,
	}
	return s
}

func (s *TaskService) CreateTask(task *models.Task) (int, error) {
	taskStatusChan := make(chan *models.TaskStatus)
	errChan := make(chan error)
	s.wp.Submit(func() {
		defer close(taskStatusChan)
		defer close(errChan)

		taskStatus, err := utils.ExecuteTask(task)
		if err != nil {
			errChan <- err
			return
		}
		taskStatusChan <- taskStatus
	})

	select {
	case taskStatus, ok := <-taskStatusChan:
		if !ok {
			return 0, errors.New("task status channel closed")
		}
		id, err := s.repo.CreateTask(task)
		if err != nil {
			return 0, err
		}
		err = s.repo.CreateTaskStatus(id, taskStatus)
		if err != nil {
			return 0, err
		}
		return id, nil
	case err, ok := <-errChan:
		if !ok {
			return 0, errors.New("error channel closed")
		}
		return 0, err
	}
}

func (s *TaskService) GetAllTasks() ([]*models.TaskFromDB, error) {
	cacheKey := "tasks_all"
	ctx := context.Background()

	cache, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var tasks []*models.TaskFromDB
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

	// caching all tasks
	err = s.redis.Set(ctx, cacheKey, jsonTasks, s.expiration).Err()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("all tasks is cached")
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

	// caching task by id
	err = s.redis.Set(ctx, cacheKey, jsonTask, s.expiration).Err()
	if err != nil {
		return nil, err
	}
	logrus.Debugf("task_%d is cached", id)
	return task, nil
}
