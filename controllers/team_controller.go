package controllers

import (
	"net/http"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
)

func ListTeams(c *gin.Context) {
	var teams []models.Team
	config.DB.Preload("Players").Find(&teams)
	c.JSON(http.StatusOK, gin.H{"data": teams})
}

func GetTeam(c *gin.Context) {
	id := c.Param("id")
	var team models.Team
	if err := config.DB.Preload("Players").Where("id = ?", id).First(&team).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Team not found", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": team})
}

func CreateTeam(c *gin.Context) {
	var body models.Team
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&body).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Message": "failed create team", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Success creating team data", "data": body})
}

func UpdateTeam(c *gin.Context) {
	id := c.Param("id")
	var team models.Team
	if err := config.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Team not found", "error": err.Error()})
		return
	}
	var body models.Team
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "JSON format incorrect", "error": err.Error()})
		return
	}
	config.DB.Model(&team).Updates(body)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updating team data", "data": team})
}

func DeleteTeam(c *gin.Context) {
	id := c.Param("id")
	var team models.Team
	if err := config.DB.First(&team, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Message": "Team not found", "error": err.Error()})
		return
	}
	config.DB.Delete(&team)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleting team"})
}
