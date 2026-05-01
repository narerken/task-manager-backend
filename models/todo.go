package models

import "time"

type Todo struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Completed   bool `gorm:"default:false"`
	Priority    int  `gorm:"default:1"`
	DueDate     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
