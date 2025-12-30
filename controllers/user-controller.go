package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/models/serializers"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"
)

var userService *services.ServiceImpl

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	userService = services.NewService(db)
	return &UserController{
		db: db,
	}
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	var user serializers.UserSerializer

	userRaw, err := userService.GetUserByID(ctx)
	httperr.HandlerError(err)

	err = copier.Copy(&user, &userRaw)
	httperr.HandlerError(err)

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (c *UserController) Me(ctx *gin.Context) {
	var user serializers.UserSerializer
	userRaw, err := userService.GetUser(ctx)
	httperr.HandlerError(err)

	err = copier.Copy(&user, &userRaw)

	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req models.UserRequest

	httperr.HandlerError(ctx.ShouldBindJSON(&req))

	err := userService.UpdateUser(ctx, &req)
	httperr.HandlerError(err)

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}
