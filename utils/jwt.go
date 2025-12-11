package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret = []byte("replace_this")

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	fmt.Println("JWT Secret in utils:", string(jwtSecret))
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	parts := strings.Split(tokenString, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		tokenString = parts[1]
	}

	fmt.Println("Validating Token:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{},
		error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
