package controllers

import (
	"fmt"
	"net/http"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/Iqbal-Maulana67/ayo-technical-test/utils"

	"github.com/gin-gonic/gin"
)

func SubmitMatchResult(c *gin.Context) {
	matchID := c.Param("id")
	var input struct {
		Goals []struct {
			PlayerID uint `json:"player_id"`
			Minute   int  `json:"minute"`
		} `json:"goals"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var match models.Match
	if err := config.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	var existsGoal []models.MatchGoal
	if err := config.DB.Where("match_id = ?", match.ID).Find(&existsGoal).Error; err == nil {
		for goal := range existsGoal {
			config.DB.Delete(&existsGoal[goal])
		}
	}

	homeScore := 0
	awayScore := 0
	for _, g := range input.Goals {
		var player models.Player
		if err := config.DB.First(&player, g.PlayerID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player not found"})
			return
		}
		if player.TeamID == match.HomeTeamID {
			homeScore++
		} else if player.TeamID == match.AwayTeamID {
			awayScore++
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "player does not belong to either team in match"})
			return
		}

		mg := models.MatchGoal{MatchID: match.ID, PlayerID: player.ID, Minute: g.Minute}
		config.DB.Create(&mg)
	}

	var matchExists models.MatchResult

	if err := config.DB.Where("match_id = ?", match.ID).First(&matchExists).Error; err == nil {
		matchExists.FinalScoreHome = homeScore
		matchExists.FinalScoreAway = awayScore
		config.DB.Save(&matchExists)
		c.JSON(http.StatusAccepted, gin.H{"message": "Match Result Updated!", "data": matchExists})
	} else {
		mr := models.MatchResult{MatchID: match.ID, FinalScoreHome: homeScore,
			FinalScoreAway: awayScore}

		if err := config.DB.Create(&mr).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed save result"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"data": mr})
	}

}

func GetMatchResult(c *gin.Context) {
	matchID := c.Param("id")

	var resultExists models.MatchResult
	if err := config.DB.Where("match_id = ?", matchID).First(&resultExists).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match result not found"})
		return
	}

	var output struct {
		HomeTeam       string `json:"home_team"`
		AwayTeam       string `json:"away_team"`
		FinalScore     string `json:"final_score"`
		Result         string `json:"result"`
		HomeTeamWins   int    `json:"home_team_wins_accumulation"`
		AwayTeamWins   int    `json:"away_team_wins_accumulation"`
		TopGoalScorers []struct {
			PlayerName string `json:"player_name"`
			TeamName   string `json:"team_name"`
			GoalCount  int    `"json:"goal_count"`
		} `json:"top_goal_scorers"`
		MatchDetails models.Match `json:"match_details"`
	}

	if err := config.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Result").Preload("Goals.Player").First(&output.MatchDetails, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}

	output.HomeTeam = output.MatchDetails.HomeTeam.Name
	output.AwayTeam = output.MatchDetails.AwayTeam.Name
	output.FinalScore = fmt.Sprintf("%d - %d", output.MatchDetails.Result.FinalScoreHome, output.MatchDetails.Result.FinalScoreAway)

	if output.MatchDetails.Result.FinalScoreHome > output.MatchDetails.Result.FinalScoreAway {
		output.Result = fmt.Sprintf("%s Wins", output.MatchDetails.HomeTeam.Name)
	} else if output.MatchDetails.Result.FinalScoreAway > output.MatchDetails.Result.FinalScoreHome {
		output.Result = fmt.Sprintf("%s Wins", output.MatchDetails.AwayTeam.Name)
	} else {
		output.Result = "Draw"
	}

	output.HomeTeamWins = utils.WinAccumulation(output.MatchDetails.HomeTeamID)
	output.AwayTeamWins = utils.WinAccumulation(output.MatchDetails.AwayTeamID)
	output.TopGoalScorers = utils.TopGoalScorers(output.MatchDetails.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": output})
}
