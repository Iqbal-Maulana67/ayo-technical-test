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
	TeamID       uint           `json:"team_id"`
	Name         string         `json:"name" gorm:"not null"`
	Height       float32        `json:"height"`
	Weight       float32        `json:"weight"`
	Position     Position       `json:"position"`
	JerseyNumber int            `json:"jersey_number" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
