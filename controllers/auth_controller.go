package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Iqbal-Maulana67/ayo-technical-test/config"
	"github.com/Iqbal-Maulana67/ayo-technical-test/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func RegisterAdmin(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	var validationMessages = map[string]string{
		"Email.required":    "Email is required",
		"Email.email":       "Email format is invalid",
		"Password.required": "Password is required",
		"Password.min":      "Password must be at least 6 characters",
	}

	if err := c.ShouldBindJSON(&input); err != nil {
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
