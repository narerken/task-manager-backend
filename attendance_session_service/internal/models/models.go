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

type AttendanceEntry struct {
	Id        int        `json:"id" db:"id"`
	SessionId int        `json:"session_id" db:"session_id"`
	StudentId int        `json:"student_id" db:"student_id"`
	Status    string     `json:"status" db:"status"`
	Comment   *string    `json:"comment,omitempty" db:"comment"`
	MarkedBy  *int       `json:"marked_by,omitempty" db:"marked_by"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
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

type CreateAttendanceSessionInput struct {
	DepartmentId int    `json:"department_id" validate:"required"`
	Date         string `json:"date" validate:"required"`
}

type AttendanceEntryUpsertInput struct {
	StudentId int     `json:"student_id" validate:"required"`
	Status    string  `json:"status" validate:"required,oneof=present absent late excused"`
	Comment   *string `json:"comment,omitempty"`
}

type BulkUpsertAttendanceEntriesInput struct {
	Entries []AttendanceEntryUpsertInput `json:"entries"`
}

type User struct {
	Id           int        `json:"id"`
	FullName     string     `json:"full_name"`
	Email        string     `json:"email"`
	DepartmentId *int       `json:"department_id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

type Department struct {
	Id        int        `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
