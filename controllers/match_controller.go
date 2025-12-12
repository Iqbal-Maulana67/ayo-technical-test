package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ListMatches(c *gin.Context) {
	var matches []models.Match
	config.DB.Preload("HomeTeam").Preload("AwayTeam").Preload("Result").Preload("Goals.Player").Find(&matches)
	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": matches})
}

func CreateMatch(c *gin.Context) {
	var body struct {
		MatchDate  string `json:"match_date" binding:"required"` // ISO datetime
		HomeTeamID uint   `json:"home_team_id" binding:"required"`
		AwayTeamID uint   `json:"away_team_id" binding:"required"`
	}

	var validationMessages = map[string]string{
		"MatchDate.required":  "Match date is required",
		"HomeTeamID.required": "Home Team ID is required",
		"AwayTeamID.required": "Away Team ID is required",
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			messages := make(map[string]string)

			for _, fe := range verr {
				key := fe.Field() + "." + fe.Tag()
				if msg, ok := validationMessages[key]; ok {
					messages[strings.ToLower(fe.Field())] = msg
				} else {
					messages[strings.ToLower(fe.Field())] = fe.Error() // fallback
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": messages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsed, err := time.Parse(time.RFC3339, body.MatchDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	if body.HomeTeamID == body.AwayTeamID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home team and away team cannot be the same"})
		return
	}
	// check if teams exist
	var homeTeam, awayTeam models.Team
	if err := config.DB.First(&homeTeam, body.HomeTeamID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home team not found"})
		return
	}

	if err := config.DB.First(&awayTeam, body.AwayTeamID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "away team not found"})
		return
	}

	match := models.Match{MatchDate: parsed, HomeTeamID: body.HomeTeamID, AwayTeamID: body.AwayTeamID}
	if err := config.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create match"})
		return
	}

	config.DB.Preload("HomeTeam").Preload("HomeTeam.Players").Preload("AwayTeam").Preload("AwayTeam.Players").Preload("MatchResult").First(&match)
	c.JSON(http.StatusCreated, gin.H{"message": "Success creating match data", "data": match})
}

func UpdateMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}
	var body struct {
		MatchDate  string `json:"match_date" binding:"required"` // ISO datetime
		HomeTeamID uint   `json:"home_team_id" binding:"required"`
		AwayTeamID uint   `json:"away_team_id" binding:"required"`
	}

	var validationMessages = map[string]string{
		"MatchDate.required":  "Match date is required",
		"HomeTeamID.required": "Home Team ID is required",
		"AwayTeamID.required": "Away Team ID is required",
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			messages := make(map[string]string)

			for _, fe := range verr {
				key := fe.Field() + "." + fe.Tag()
				if msg, ok := validationMessages[key]; ok {
					messages[strings.ToLower(fe.Field())] = msg
				} else {
					messages[strings.ToLower(fe.Field())] = fe.Error() // fallback
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": messages})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	parsed, err := time.Parse(time.RFC3339, body.MatchDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	if body.HomeTeamID == body.AwayTeamID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home team and away team cannot be the same"})
		return
	}

	// check if teams exist
	var homeTeam, awayTeam models.Team
	if err := config.DB.First(&homeTeam, body.HomeTeamID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "home team not found"})
		return
	}

	if err := config.DB.First(&awayTeam, body.AwayTeamID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "away team not found"})
		return
	}

	match.MatchDate = parsed
	match.HomeTeamID = body.HomeTeamID
	match.AwayTeamID = body.AwayTeamID
	config.DB.Save(&match)
	config.DB.Preload("HomeTeam").Preload("HomeTeam.Players").Preload("AwayTeam").Preload("AwayTeam.Players").Preload("MatchResult").First(&match)
	c.JSON(http.StatusOK, gin.H{"message": "Success updating match data", "data": match})
}

func DeleteMatch(c *gin.Context) {
	id := c.Param("id")
	var match models.Match
	if err := config.DB.First(&match, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}
	config.DB.Delete(&match)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleting match"})
}
