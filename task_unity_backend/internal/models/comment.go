package models

import "time"

type Comment struct {
	Id        int        `json:"id" db:"id"`
	Comment   string     `json:"comment" db:"comment"`
	TaskId    int        `json:"task_id" db:"task_id"`
	CreatorId int        `json:"creator_id" db:"user_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}
