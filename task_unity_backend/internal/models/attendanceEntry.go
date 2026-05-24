package models

import "time"

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
