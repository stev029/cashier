package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/models/serializers"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"
)

var ItemService *services.ItemService

type ItemControllerImpl struct {
	db *gorm.DB
}

func NewItemController(db *gorm.DB) *ItemControllerImpl {
	ItemService = services.NewItemService(db)
	return &ItemControllerImpl{db: db}
}

func (c *ItemControllerImpl) GetItems(ctx *gin.Context) {
	itemsRaw, err := ItemService.GetItems(ctx)
	var items []serializers.ItemSerializer
	if err != nil {
		log.Printf("Error while getting items: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := copier.Copy(&items, &itemsRaw); err != nil {
		log.Printf("Error while copying items: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(200, gin.H{"data": items})
}

func (c *ItemControllerImpl) GetItemByID(ctx *gin.Context) {
	var item serializers.ItemSerializer
	itemRaw, err := ItemService.GetItemByID(ctx)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(404, gin.H{"error": "Item not found"})
			return
		}
		log.Printf("Error while getting item: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	if err := copier.Copy(&item, &itemRaw); err != nil {
		log.Printf("Error while copying item: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(200, gin.H{"data": item})
}

func (c *ItemControllerImpl) CreateItem(ctx *gin.Context) {
	var ItemRequest models.ItemRequest

	if err := ctx.ShouldBindJSON(&ItemRequest); err != nil {
		log.Printf("Error while binding JSON: %v", err)
		ctx.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	_, err := ItemService.CreateItem(ctx, ItemRequest)
	if err != nil {
		log.Printf("Error while creating item: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(201, gin.H{"status": "success"})
}

func (c *ItemControllerImpl) UpdateItem(ctx *gin.Context) {
	var ItemRequest models.ItemRequest
	if err := ctx.ShouldBindJSON(&ItemRequest); err != nil {
		log.Printf("Error while binding JSON: %v", err)
		return
	}

	_, err := ItemService.UpdateItem(ctx, ItemRequest)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(404, gin.H{"error": "Item not found"})
			return
		}
		log.Printf("Error while updating item: %v", err)
		ctx.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ctx.JSON(200, gin.H{"status": "success"})
}

func (c *ItemControllerImpl) DeleteItem(ctx *gin.Context) {
	ItemService.DeleteItem(ctx)
	ctx.JSON(200, gin.H{"status": "success"})
}
