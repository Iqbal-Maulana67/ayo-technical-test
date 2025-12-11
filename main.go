package main

import (
	"log"
	"os"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/Iqbal-Maulana67/ayo-technical-test/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	dsn := os.Getenv("DATABASE_URL")
	db := config.ConnectDB(dsn)
	if err := db.AutoMigrate(&models.Team{}, &models.Player{}, &models.Match{},
		&models.MatchResult{}, &models.MatchGoal{}, &models.UserAdmin{}); err != nil {
		log.Fatal("migrate failed", err)
	}
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
