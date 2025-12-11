package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("AYO_TECHNICAL_TEST_SECRET")

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
}

// ===============================
// REGISTER ADMIN
// ===============================
func RegisterAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check existing admin
	var existing models.UserAdmin
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Email Registered"})
		return
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	admin := models.UserAdmin{
		Email:    input.Email,
		Password: string(hashed),
	}

	config.DB.Create(&admin)
	c.JSON(http.StatusCreated, gin.H{"message": "data successfully created", "data": admin})
}

// ===============================
// LOGIN ADMIN
// ===============================
func LoginAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get admin by email
	var admin models.UserAdmin
	if err := config.DB.Where("email = ?", input.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Acccount not found"})
		return
	}

	// Check password
	if bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password incorrect"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": input.Email,
			"exp":   time.Now().Add(24 * time.Hour).Unix(),
		})

	signedToken, _ := token.SignedString(jwtSecret)
	fmt.Println("JWT Secret:", jwtSecret)
	fmt.Println("Token:", signedToken)

	c.JSON(http.StatusAccepted, gin.H{"message": "Login success", "token": signedToken})
}
