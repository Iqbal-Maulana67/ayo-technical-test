package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchGoal struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	MatchID   uint           `json:"match_id"`
	PlayerID  uint           `json:"player_id"`
	Minute    int            `json:"minute"`
	Player    Player         `gorm:"foreignKey:PlayerID" json:"player"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
