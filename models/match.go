package models

import (
	"time"

	"gorm.io/gorm"
)

type Match struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	MatchDate  time.Time      `json:"match_date"` // store date+time
	HomeTeamID uint           `json:"home_team_id"`
	AwayTeamID uint           `json:"away_team_id"`
	HomeTeam   Team           `gorm:"foreignKey:HomeTeamID" json:"home_team"`
	AwayTeam   Team           `gorm:"foreignKey:AwayTeamID" json:"away_team"`
	Result     MatchResult    `json:"result" gorm:"constraint:OnDelete:SET NULL;"`
	Goals      []MatchGoal    `json:"goals"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
