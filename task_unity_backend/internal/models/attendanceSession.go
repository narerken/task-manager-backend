package models

import "time"

type AttendanceSession struct {
	Id           int        `json:"id" db:"id"`
	DepartmentId int        `json:"department_id" db:"department_id"`
	Date         time.Time  `json:"date" db:"date"`
	State        string     `json:"state" db:"state"`
	CreatedBy    *int       `json:"created_by,omitempty" db:"created_by"`
	UpdatedBy    *int       `json:"updated_by,omitempty" db:"updated_by"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type AttendanceSessionEntryView struct {
	StudentId   int     `json:"student_id"`
	StudentName string  `json:"student_name"`
	Status      *string `json:"status,omitempty"`
	Comment     *string `json:"comment,omitempty"`
	MarkedBy    *int    `json:"marked_by,omitempty"`
}

type AttendanceSessionDetails struct {
	Id           int                          `json:"id"`
	DepartmentId int                          `json:"department_id"`
	Date         string                       `json:"date"`
	State        string                       `json:"state"`
	CreatedBy    *int                         `json:"created_by,omitempty"`
	UpdatedBy    *int                         `json:"updated_by,omitempty"`
	Entries      []AttendanceSessionEntryView `json:"entries"`
}
