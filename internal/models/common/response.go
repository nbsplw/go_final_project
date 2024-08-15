package common

import "main/internal/models/tasks"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Id    int          `json:"id,omitempty"`
	Tasks []tasks.Task `json:"tasks,omitempty"`
}
