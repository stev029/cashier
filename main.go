package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stev029/cashier/controllers"
	"github.com/stev029/cashier/etc/database"
	_ "github.com/stev029/cashier/etc/database/autoload"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	db = database.DB
	err := database.InitModel()
	if err != nil {
		if err == gorm.ErrInvalidDB {
			log.Fatalf("DB not initialized: %v", err)
			return
		}

		log.Fatalf("Error while init model: %v", err)
		return
	}
}

func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	log.Print("starting server...")
	router.GET("/ping", handler)
	controllers.Controller(router, db)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func handler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Pong!"})
}
