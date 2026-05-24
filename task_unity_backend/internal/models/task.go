package models

import "time"

type Task struct {
	Id           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Description  string     `json:"description" db:"description"`
	Deadline     time.Time  `json:"deadline" db:"deadline"`
	DepartmentId int        `json:"department_id" db:"department_id"`
	CreatorId    int        `json:"creator_id" db:"creator_id"`
	AssigneeId   int        `json:"assignee_id" db:"assignee_id"`
	Status       string     `json:"status" validate:"oneof=todo in_progress done"  db:"status"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
