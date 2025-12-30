package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/models/serializers"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"
)

var ItemService *services.ServiceImpl

type ItemControllerImpl struct {
	db *gorm.DB
}

func NewItemController(db *gorm.DB) *ItemControllerImpl {
	ItemService = services.NewService(db)
	return &ItemControllerImpl{db: db}
}

func (c *ItemControllerImpl) GetItems(ctx *gin.Context) {
	var items []serializers.ItemSerializer
	itemsRaw, err := ItemService.GetItems(ctx)
	httperr.HandlerError(err)

	err = copier.Copy(&items, &itemsRaw)
	httperr.HandlerError(err)

	ctx.JSON(200, gin.H{"data": items})
}

func (c *ItemControllerImpl) GetItemByID(ctx *gin.Context) {
	var item serializers.ItemSerializer
	itemRaw, err := ItemService.GetItemByID(ctx)
	httperr.HandlerError(err)

	err = copier.Copy(&item, &itemRaw)
	httperr.HandlerError(err)

	ctx.JSON(200, gin.H{"data": item})
}

func (c *ItemControllerImpl) CreateItem(ctx *gin.Context) {
	var ItemRequest models.ItemRequest

	httperr.HandlerError(ctx.ShouldBindJSON(&ItemRequest))

	_, err := ItemService.CreateItem(ctx, ItemRequest)
	httperr.HandlerError(err)

	ctx.JSON(201, gin.H{"status": "success"})
}

func (c *ItemControllerImpl) UpdateItem(ctx *gin.Context) {
	var ItemRequest models.ItemRequest
	httperr.HandlerError(ctx.ShouldBindJSON(&ItemRequest))

	_, err := ItemService.UpdateItem(ctx, ItemRequest)
	httperr.HandlerError(err)

	ctx.JSON(200, gin.H{"status": "success"})
}

func (c *ItemControllerImpl) DeleteItem(ctx *gin.Context) {
	ItemService.DeleteItem(ctx)
	ctx.JSON(200, gin.H{"status": "success"})
}
