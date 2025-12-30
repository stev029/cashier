package services

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/etc/database"
	"github.com/stev029/cashier/etc/database/model"
	"github.com/stev029/cashier/etc/utils"
	"github.com/stev029/cashier/models"

	"github.com/midtrans/midtrans-go/coreapi"
)

type TransactionService interface {
	CreateTransaction(transaction model.Transaction) (*coreapi.ChargeResponse, error)
}

func (s *ServiceImpl) CreateTransaction(ctx *gin.Context, transaction models.TransactionRequest) (*utils.DuitkuResponseCharge, error) {
	order_id := fmt.Sprintf("order-%d", time.Now().Unix())
	chargeReq := utils.DuitkuRequestCharge{
		MerchantOrderId: order_id,
		ProductDetails:  fmt.Sprintf("Payment-Invoice-%s", order_id),
		Email:           "abyy1144@gmail.com",
		PaymentMethod:   utils.GUDANGVQR,
		CustomerVaName:  transaction.CustomerName,
		PaymentAmount:   int(transaction.Price),
	}

	resp, _, err := utils.DuitkuClient.PaymentService.Charge(ctx, chargeReq)
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
		PaymentCode: resp.QrString, // For QRIS, the payment code is in the URL
	}

	if err := db.Create(&newTransaction).Error; err != nil {
		return nil, err
	}

	return &resp, nil
}
