package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/etc/database"
	"github.com/stev029/cashier/etc/database/model"
	"github.com/stev029/cashier/etc/utils"
	"github.com/stev029/cashier/models"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type TransactionService interface {
	CreateTransaction(transaction model.Transaction) (*coreapi.ChargeResponse, error)
}

func (s *ServiceImpl) CreateTransaction(ctx *gin.Context, transaction models.TransactionRequest) (*coreapi.ChargeResponse, error) {
	order_id := fmt.Sprintf("order-%d", time.Now().Unix())
	chargeReq := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  order_id,
			GrossAmt: int64(transaction.Price),
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "item-1",
				Name:  "Item 1",
				Price: int64(transaction.Price),
				Qty:   1,
			},
		},
		Qris: &coreapi.QrisDetails{
			Acquirer: "airpay shopee",
		},
	}

	resp, err := utils.MidtransClient.ChargeTransaction(chargeReq)
	if err != nil {
		return nil, err
	}

	// Save transaction to database
	db := database.DB
	newTransaction := model.Transaction{
		UserID:      transaction.UserID,
		Amount:      transaction.Amount,
		OrderID:     order_id,
		Status:      resp.StatusCode,
		PaymentType: transaction.PaymentType,
		PaymentCode: resp.QRString, // For QRIS, the payment code is in the URL
	}

	if err := db.Create(&newTransaction).Error; err != nil {
		return nil, err
	}

	return resp, nil
}
