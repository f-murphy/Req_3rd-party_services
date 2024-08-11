package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"req3rdPartyServices/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TaskRepositoryInterface interface {
	CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error)
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskRepository struct {
	db         *sqlx.DB
	redis      *redis.Client
	expiration time.Duration
}

func NewTaskRepository(db *sqlx.DB, redis *redis.Client, expiration time.Duration) *TaskRepository {
	return &TaskRepository{db: db, redis: redis, expiration: expiration}
}

func (r *TaskRepository) CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error) {
	var id int
	queryCreateTask := `INSERT INTO Tasks (Method, Url, Headers, Body) VALUES ($1, $2, $3, $4) RETURNING Id`
	err := r.db.QueryRow(queryCreateTask, task.Method, task.Url, task.HeadersJSON, task.BodyJSON).Scan(&id)
	if err != nil {
		return 0, err
	}

	queryTaskStatus := `INSERT INTO TaskStatus (Status, HttpStatusCode, Length) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(queryTaskStatus, taskStatus.Status, taskStatus.HttpStatusCode, taskStatus.Length)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (r *TaskRepository) GetAllTasks() ([]*models.TaskFromDB, error) {
	cacheKey := "tasks_all"
	ctx := context.Background()

	cache, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var tasks []*models.TaskFromDB
		err = json.Unmarshal([]byte(cache), &tasks)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	}

	tasks := []*models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
	`
	err = r.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}

	jsonTasks, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}

	err = r.redis.Set(ctx, cacheKey, jsonTasks, r.expiration).Err()
	if err != nil {
		return nil, err
	}
	logrus.Infof("tasks_all is cached")

	return tasks, nil
}

func (r *TaskRepository) GetTaskById(id int) (*models.TaskFromDB, error) {
	cacheKey := fmt.Sprintf("task_%d", id)
	ctx := context.Background()

	cache, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var task *models.TaskFromDB
		err = json.Unmarshal([]byte(cache), &task)
		if err != nil {
			return nil, err
		}
		return task, nil
	}

	task := &models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
		WHERE tasks.id = $1
	`
	err = r.db.Get(task, query, id)
	if err != nil {
		return nil, err
	}

	jsonTask, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	err = r.redis.Set(ctx, cacheKey, jsonTask, r.expiration).Err()
	if err != nil {
		return nil, err
	}
	logrus.Infof("task_%d is cached", id)

	return task, nil
}
