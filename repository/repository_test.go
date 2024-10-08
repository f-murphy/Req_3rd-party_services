package repository

import (
	"reflect"
	"req3rdPartyServices/models"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestTaskRepository_CreateTask(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres:qwerty@localhost:5436/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	repo := NewTaskRepository(db)

	type args struct {
		task       *models.Task
		taskStatus *models.TaskStatus
	}
	tests := []struct {
		name    string
		r       *TaskRepository
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "CreateTaskSuccess",
			r:    repo,
			args: struct {
				task       *models.Task
				taskStatus *models.TaskStatus
			}{
				task: &models.Task{
					Method: "GET",
					Url:    "google.com",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]string{
						"test": "test",
					},
				},
				taskStatus: &models.TaskStatus{},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "CreateTaskFailed",
			r:    repo,
			args: struct {
				task       *models.Task
				taskStatus *models.TaskStatus
			}{
				task: &models.Task{
					Method: "GET",
					Url:    "gggoogle.com",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]string{
						"test": "test",
					},
				},
				taskStatus: &models.TaskStatus{},
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.CreateTask(tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.CreateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskRepository.CreateTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepository_GetTaskById(t *testing.T) {
	db, err := sqlx.Connect("postgres", "postgres:qwerty@localhost:5436/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	repo := NewTaskRepository(db)

	type args struct {
		id int
	}
	tests := []struct {
		name    string
		r       *TaskRepository
		args    args
		want    *models.TaskFromDB
		wantErr bool
	}{
		{
			name: "GetTaskByIdSuccess",
			r:    repo,
			args: struct {
				id int
			}{
				id: 1,
			},
			want: &models.TaskFromDB{
				Id:             1,
				Method:         "GET",
				Url:            "https://google.com",
				Headers:        "",
				Body:           "",
				Status:         "200",
				HttpStatusCode: "200 OK",
				Length:         "20",
			},
		},
		{
			name: "GetTaskByIdError",
			r:    repo,
			args: struct {
				id int
			}{
				id: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetTaskById(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.GetTaskById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepository.GetTaskById() = %v, want %v", got, tt.want)
			}
		})
	}
}
