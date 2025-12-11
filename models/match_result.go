package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchResult struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	MatchID        uint           `json:"match_id" gorm:"uniqueIndex"`
	FinalScoreHome int            `json:"final_score_home"`
	FinalScoreAway int            `json:"final_score_away"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
