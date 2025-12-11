package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `json:"name" gorm:"not null"`
	LogoURL     string         `json:"logo_url"`
	FoundedYear int            `json:"founded_year"`
	HomeAddress string         `json:"home_address"`
	HomeCity    string         `json:"home_city"`
	Players     []Player       `json:"players" gorm:"constraint:OnDelete:SET NULL;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
