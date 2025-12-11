package models

import (
	"time"

	"gorm.io/gorm"
)

type UserAdmin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
