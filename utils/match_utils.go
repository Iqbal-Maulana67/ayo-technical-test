package utils

import (
	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
)

func WinAccumulation(teamID uint) int {

	var matchList []models.Match
	var WinsTotal int

	if err := config.DB.Where("home_team_id = ? OR away_team_id = ?", teamID, teamID).Find(&matchList).Error; err != nil {
		return 0
	}

	for _, match := range matchList {
		var matchResult models.MatchResult
		if err := config.DB.Where("match_id = ?", match.ID).First(&matchResult).Error; err != nil {
			continue
		}

		homeGoals := matchResult.FinalScoreHome
		awayGoals := matchResult.FinalScoreAway

		if teamID == match.HomeTeamID {
			if homeGoals > awayGoals {
				WinsTotal++
			}
		} else if teamID == match.AwayTeamID {
			if awayGoals > homeGoals {
				WinsTotal++
			}
		}
	}
	return WinsTotal
}

func TopGoalScorers(matchID uint) []struct {
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	GoalCount  int    `"json:"goal_count"`
} {
	var results []struct {
		PlayerName string `json:"player_name"`
		TeamName   string `json:"team_name"`
		GoalCount  int    `"json:"goal_count"`
	}

	config.DB.Table("match_goals").Select("players.name as player_name, teams.name as team_name, COUNT(match_goals.id) as goal_count").
		Where("match_id = ? AND match_goals.deleted_at IS NULL", matchID).
		Joins("JOIN players ON match_goals.player_id = players.id").
		Joins("JOIN teams ON players.team_id = teams.id").
		Group("players.id, teams.name").
		Order("goal_count DESC").
		Scan(&results)
	return results
}
