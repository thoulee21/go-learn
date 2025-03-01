package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status" binding:"oneof=pending completed canceled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
