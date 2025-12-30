package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	httperr "github.com/stev029/cashier/http-errors"
)

func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Error: %v", err.(error).Error())
				if err, ok := err.(*httperr.HttpExceptionJSONImpl); ok {
					c.AbortWithStatusJSON(http.StatusNotFound, err.Message)
					return
				}

				// Jika error lain, kirim 500 Internal Server Error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "internal server error",
				})
			}
		}()
		c.Next()
	}
}
