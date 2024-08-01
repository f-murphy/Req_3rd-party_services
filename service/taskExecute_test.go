package service

import (
	"reflect"
	"req3rdPartyServices/models"
	"testing"
)

func TestExecuteTask(t *testing.T) {
	tests := []struct {
		name string
		task *models.Task
		wantStatus *models.TaskStatus
		wantErr bool
	} {
		{
			name: "success",
			task: &models.Task{
				Method: "GET",
				Url: "",
			},
			wantStatus: &models.TaskStatus{
				Status: "200 OK",
				HttpStatusCode: "200",
				Length: "",
			},
			wantErr: false,
		},
		{
			name: "success",
			task: &models.Task{
				Method: "GET",
				Url: "",
			},
			wantStatus: &models.TaskStatus{
				Status: "200 OK",
				HttpStatusCode: "200",
				Length: "",
			},
			wantErr: false,
		},
		{
			name: "success",
			task: &models.Task{
				Method: "GET",
				Url: "",
			},
			wantStatus: &models.TaskStatus{
				Status: "200 OK",
				HttpStatusCode: "200",
				Length: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, err := ExecuteTask(tt.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStatus, tt.wantStatus) {
				t.Errorf("ExecuteTask() = %v, want %v", gotStatus, tt.wantStatus)
			}
		})
	}
}