package repository

import (
	"req3rdPartyServices/models"

	"github.com/jmoiron/sqlx"
)

type TaskRepositoryInterface interface {
	CreateTask(task *models.Task) (int, error)
	CreateTaskStatus(taskID int, taskStatus *models.TaskStatus) error
	GetAllTasks() ([]*models.TaskFromDB, error)
	GetTaskById(id int) (*models.TaskFromDB, error)
}

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *models.Task) (int, error) {
	var id int
	queryCreateTask := `INSERT INTO Tasks (Method, Url, Headers, Body) VALUES ($1, $2, $3, $4) RETURNING Id`
	err := r.db.QueryRow(queryCreateTask, task.Method, task.Url, task.HeadersJSON, task.BodyJSON).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TaskRepository) CreateTaskStatus(taskID int, taskStatus *models.TaskStatus) error {
	queryTaskStatus := `INSERT INTO TaskStatus (Id, Status, HttpStatusCode, Length) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(queryTaskStatus, taskID, taskStatus.Status, taskStatus.HttpStatusCode, taskStatus.Length)
	return err
}

func (r *TaskRepository) GetAllTasks() ([]*models.TaskFromDB, error) {
	tasks := []*models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
	`
	err := r.db.Select(&tasks, query)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) GetTaskById(id int) (*models.TaskFromDB, error) {
	task := &models.TaskFromDB{}
	query := `
		SELECT * FROM Tasks
		INNER JOIN TaskStatus ON Tasks.id = TaskStatus.id
		WHERE tasks.id = $1
	`
	err := r.db.Get(task, query, id)
	if err != nil {
		return nil, err
	}
	return task, nil
}
