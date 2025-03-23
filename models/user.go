package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	UserName     string    `json:"user_name" gorm:"unique;not null"`
	HashPassword string    `json:"hash_password" gorm:"not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type IUserService interface {
	Create(newUser *User) (*User, error)
	Delete(id uint) error
	Update(id uint, updatedUser *User) (*User, error)
	GetAll() (*[]User, error)
	GetByID(id uint) (*User, error)
	GetOneByMap(userMap map[string]any) (*User, error)
}
