package models

import (
	"time"

	"gorm.io/gorm"
)

type Position string

const (
	Penyerang     Position = "penyerang"
	Gelandang     Position = "gelandang"
	Bertahan      Position = "bertahan"
	PenjagaGawang Position = "penjaga_gawang"
)

type Player struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	TeamID       uint           `json:"team_id" binding:"required"`
	Name         string         `json:"name" gorm:"not null" binding:"required"`
	Height       float32        `json:"height" binding:"required"`
	Weight       float32        `json:"weight" binding:"required"`
	Position     Position       `json:"position" binding:"required"`
	JerseyNumber int            `json:"jersey_number" gorm:"not null" binding:"required"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
