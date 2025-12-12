package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ListPlayers(c *gin.Context) {
	teamID := c.Query("team_id")
	var players []models.Player
	if teamID != "" {
		config.DB.Where("team_id = ?", teamID).Find(&players)
	} else {
		config.DB.Find(&players)
	}
	c.JSON(http.StatusOK, gin.H{"data": players})
}

func CreatePlayer(c *gin.Context) {
	var body models.Player

	var validationMessages = map[string]string{
		"TeamID.required":       "Match date is required",
		"Name.required":         "Name is required",
		"Height.required":       "Height is required",
		"Weight.required":       "Weight is required",
		"Position.required":     "Position is required",
		"JerseyNumber.required": "Jersey Number is required",
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "JSON format error", "error": err.Error()})
		return
	}

	var exists models.Player
	var teamExists models.Team

	if err := config.DB.Where("ID = ? AND deleted_at IS NULL", body.TeamID).First(&teamExists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Team not found"})
		return
	}

	// check unique jersey number per team
	if err := config.DB.Where("team_id = ? AND jersey_number = ?", body.TeamID, body.JerseyNumber).First(&exists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Jersey number already exists in the team"})
		return
	}

	if err := config.DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed create player", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Successfully creating player", "data": body})
}

func UpdatePlayer(c *gin.Context) {
	id := c.Param("id")
	var player models.Player
	if err := config.DB.First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Player not found"})
		return
	}
	var body models.Player

	var validationMessages = map[string]string{
		"TeamID.required":       "Match date is required",
		"Name.required":         "Name is required",
		"Height.required":       "Height is required",
		"Weight.required":       "Weight is required",
		"Position.required":     "Position is required",
		"JerseyNumber.required": "Jersey Number is required",
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

		c.JSON(http.StatusBadRequest, gin.H{"message": "JSON format error", "error": err.Error()})
		return
	}

	var teamExists models.Team

	if err := config.DB.Where("ID = ? AND deleted_at IS NULL", body.TeamID).First(&teamExists).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Team not found"})
		return
	}

	var exists models.Player

	if err := config.DB.Where("team_id = ? AND jersey_number = ?", body.TeamID, body.JerseyNumber).First(&exists).Error; err == nil {
		if exists.ID != player.ID {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Jersey number already exists in the team"})
			return
		}
	}

	config.DB.Model(&player).Updates(body)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updating player data", "data": player})
}

func DeletePlayer(c *gin.Context) {
	id := c.Param("id")
	var player models.Player
	if err := config.DB.First(&player, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Player not found"})
		return
	}
	config.DB.Delete(&player)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleting player"})
}
