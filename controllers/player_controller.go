package controllers

import (
	"net/http"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
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
	if err := c.ShouldBindJSON(&body); err != nil {
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
	if err := c.ShouldBindJSON(&body); err != nil {
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
