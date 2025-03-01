package models

import (
	"time"
)

type ChatMessage struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	SessionID string    `json:"session_id" gorm:"index"`
	Role      string    `json:"role" binding:"required,oneof=user assistant system"` // user, assistant, system
	Content   string    `json:"content" binding:"required"`
}

type ChatRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message" binding:"required"`
}

type ChatResponse struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}
