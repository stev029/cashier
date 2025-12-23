package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/etc/database"
	"github.com/stev029/cashier/etc/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// CHEACK REDIS: Are have token in blacklist?
		badToken, _ := database.RedisClient.Get(c, tokenString).Result()
		if badToken == "blacklisted" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Verification JWT
		token, claims, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			log.Printf("Error verifying token: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
