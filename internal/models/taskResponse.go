package models

type TaskResponse struct {
	TaskID     int
	Task       Task
	TaskStatus any // честно, тут я застрял.. возможно решение найду позже и залью в свой репо, но пока any
}
