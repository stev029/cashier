package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	httperr "github.com/stev029/cashier/http-errors"
)

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
		"user_id": userID,
	}

	log.Println(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func VerifyToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, httperr.NewHttpExceptionJSON(500, gin.H{"error": fmt.Sprintf("method signing not valid: %v", token.Header["alg"])})
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, nil, httperr.NewHttpExceptionJSON(401, gin.H{"error": "Invalid Token"})
		}
		return nil, nil, httperr.InternalServerError
	}

	if !token.Valid {
		return nil, nil, httperr.NewHttpExceptionJSON(401, gin.H{"error": "Invalid Token"})
	}

	return token, claims, nil
}
