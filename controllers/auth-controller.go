package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"
)

var AuthService *services.ServiceImpl

type AuthControllerImpl struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthControllerImpl {
	AuthService = services.NewServiceImpl(db)
	return &AuthControllerImpl{db: db}
}

func (c *AuthControllerImpl) Register(ctx *gin.Context) {
	var RegisterRequest models.RegisterRequest

	if err := ctx.ShouldBindJSON(&RegisterRequest); err != nil {
		log.Println("Error while binding JSON:", err)
		ctx.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	token, err := AuthService.Register(ctx, RegisterRequest)
	if err != nil {
		log.Printf("Error while registering user: %v", err)
		return
	}

	ctx.JSON(201, gin.H{"status": "success", "token": token})
}

func (c *AuthControllerImpl) Login(ctx *gin.Context) {
	var LoginRequest models.LoginRequest

	if err := ctx.ShouldBindJSON(&LoginRequest); err != nil {
		log.Println("Error while binding JSON:", err)
		ctx.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	token, err := AuthService.Login(ctx, LoginRequest)
	if err != nil {
		log.Printf("Error while logging in user: %v", err)
		return
	}

	ctx.JSON(200, gin.H{"status": "success", "token": token})
}

func (c *AuthControllerImpl) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	if err := AuthService.LogoutService(ctx, tokenString); err != nil {
		log.Printf("Error while logging out user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(200, gin.H{"status": "success"})
}
