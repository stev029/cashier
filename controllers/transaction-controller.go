package controllers

import (
	"log"
	"net/http"

	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
	"github.com/stev029/cashier/services"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var transactionService *services.ServiceImpl

type TransactionController struct {
	db *gorm.DB
}

func NewTransactionController(db *gorm.DB) *TransactionController {
	transactionService = services.NewService(db)
	return &TransactionController{db: db}
}

func (c *TransactionController) CreateTransaction(ctx *gin.Context) {
	var transaction models.TransactionRequest
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		panic(err)
	}

	resp, err := transactionService.CreateTransaction(ctx, transaction)
	httperr.HandlerError(err)

	ctx.JSON(http.StatusCreated, resp)
}

func (c *TransactionController) WebhookTransaction(ctx *gin.Context) {
	bodyBytes, err := ctx.GetRawData()
	if err != nil {
		panic(httperr.BadRequest)
	}

	body := string(bodyBytes)
	log.Println(body)
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}
