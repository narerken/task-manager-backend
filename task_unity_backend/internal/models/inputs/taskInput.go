package inputs

import "time"

type AddTaskInput struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Deadline     time.Time `json:"deadline"`
	DepartmentId int       `json:"department_id"`
	AssigneeId   int       `json:"assignee_id"`
	Status       string    `json:"status" validate:"oneof=todo in_progress done"`
}

type UpdateTaskInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	AssigneeId  *int       `json:"assignee_id"`
	Status      *string    `json:"status" validate:"oneof=todo in_progress done"`
}
