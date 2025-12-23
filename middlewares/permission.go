package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/etc/database"
	"github.com/stev029/cashier/etc/database/model"
)

func RequirePermission(requiredPerm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Getting user from database with groups and permissions
		var user model.User
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// Query Join for get all of permission user
		err := database.DB.Preload("Groups.Permissions").First(&user, userID).Error
		if err != nil {
			c.AbortWithStatusJSON(403, gin.H{"error": "denied permission"})
			return
		}

		// Check are have permission match
		hasPerm := false
		for _, group := range user.Groups {
			if group.Name == "fullaccess" {
				hasPerm = true
				break
			}

			for _, p := range group.Permissions {
				if p.Name == requiredPerm {
					hasPerm = true
					break
				}
			}
		}

		if !hasPerm {
			c.AbortWithStatusJSON(403, gin.H{"error": "denied permission"})
			return
		}
		c.Next()
	}
}

func GroupPermission(requiredPerm string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Getting user from database with groups and permissions
		var user model.User
		userID, exists := ctx.Get("user_id")
		if !exists {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		err := database.DB.Preload("Groups").First(&user, userID).Error
		if err != nil {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "denied permission"})
			return
		}

		hasPerm := false
		for _, group := range user.Groups {
			if group.Name == "fullaccess" || group.Name == requiredPerm {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "denied permission"})
			return
		}
		ctx.Next()
	}
}
