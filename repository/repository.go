package repository

import (
	"req3rdPartyServices/models"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type TaskRepositoryInterface interface {
	CreateTask(task *models.Task, taskStatus *models.TaskStatus) (int, error)
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskRepository struct {
	db       *sqlx.DB
	redis    *redis.Client
	cacheTTL time.Duration
}

func NewTaskRepository(db *sqlx.DB, redis *redis.Client, cacheTTL time.Duration) *TaskRepository {
	return &TaskRepository{db: db, redis: redis, cacheTTL: cacheTTL}
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
	tasks := []*models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
	`
	err := r.db.Select(&tasks, query)
	return tasks, err
}

func (r *TaskRepository) GetTaskById(id int) (*models.TaskFromDB, error) {
	task := &models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
		WHERE tasks.id = $1
	`

	err := r.db.Get(task, query, id)
	return task, err
}
