package models

import "time"

type Attendance struct {
	Id           int        `json:"id" db:"id"`
	Date         time.Time  `json:"date" db:"date"`
	UserId       int        `json:"user_id" db:"user_id"`
	DepartmentId *int       `json:"department_id,omitempty" db:"department_id"`
	Status       *string    `json:"status,omitempty" validate:"oneof=present absent excused" db:"status"`
	Comment      *string    `json:"comment,omitempty" db:"comment"`
	MarkedBy     *int       `json:"marked_by,omitempty" db:"marked_by"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
