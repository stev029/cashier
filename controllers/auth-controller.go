package controllers

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"
)

var AuthService *services.ServiceImpl

type AuthControllerImpl struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthControllerImpl {
	AuthService = services.NewService(db)
	return &AuthControllerImpl{db: db}
}

func (c *AuthControllerImpl) Register(ctx *gin.Context) {
	var RegisterRequest models.RegisterRequest

	if err := ctx.ShouldBindJSON(&RegisterRequest); err != nil {
		log.Println("Error while binding JSON:", err)
		panic(httperr.BadRequest)
	}

	token, err := AuthService.Register(ctx, RegisterRequest)
	httperr.HandlerError(err)

	ctx.JSON(201, gin.H{"status": "success", "token": token})
}

func (c *AuthControllerImpl) Login(ctx *gin.Context) {
	var LoginRequest models.LoginRequest

	if err := ctx.ShouldBindJSON(&LoginRequest); err != nil {
		log.Println("Error while binding JSON:", err)
		panic(httperr.BadRequest)
	}

	token, err := AuthService.Login(ctx, LoginRequest)
	httperr.HandlerError(err)

	ctx.JSON(200, gin.H{"status": "success", "token": token})
}

func (c *AuthControllerImpl) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	httperr.HandlerError(AuthService.LogoutService(ctx, tokenString))

	ctx.JSON(200, gin.H{"status": "success"})
}
