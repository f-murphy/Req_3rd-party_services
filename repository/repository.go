package repository

import (
	"encoding/json"
	"fmt"
	"req3rdPartyServices/models"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *models.Task, taskStatus *models.TaskStatus) error {
	jsonHeaders, err := json.Marshal(task.Headers)
	if err != nil {
		fmt.Println(err)
	}
	
	jsonBody, err := json.Marshal(task.Body)
	if err != nil {
		fmt.Println(err)
	}

	query := `INSERT INTO Tasks (Method, Url, Headers, Body, Status, HttpStatusCode, Length) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = r.db.Exec(query, task.Method, task.Url, string(jsonHeaders), string(jsonBody), taskStatus.Status, taskStatus.HttpStatusCode, taskStatus.Length)
	return err
}

func (r *TaskRepository) GetTask(id int) (*models.TaskFromDB, error) {
	task := &models.TaskFromDB{}
	query := `SELECT * FROM Tasks WHERE id = $1`
	
	err := r.db.Get(task, query, id)
	return task, err
}