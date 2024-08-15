package repository

import (
	"reflect"
	"req3rdPartyServices/models"
	"testing"
)

func TestTaskRepository_CreateTask(t *testing.T) {
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.CreateTask(tt.args.task, tt.args.taskStatus)
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

func TestTaskRepository_GetAllTasks(t *testing.T) {
	tests := []struct {
		name    string
		r       *TaskRepository
		want    []*models.TaskFromDB
		wantErr bool
	}{
		//  Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.GetAllTasks()
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.GetAllTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskRepository.GetAllTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepository_GetTaskById(t *testing.T) {
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
		//  Add test cases.
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
