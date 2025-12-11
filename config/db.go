package config

import (
	"log"
	"os"

	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) *gorm.DB {
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	DB = db
	return DB
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Team{},
		&models.Player{},
		&models.Match{},
		&models.MatchResult{},
		&models.MatchGoal{},
	)
}
