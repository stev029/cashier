package models

type TransactionRequest struct {
	UserID       uint    `json:"user_id"`
	CustomerName string  `json:"customer_name" validate:"required"`
	Amount       int     `json:"amount"`
	Price        float64 `json:"price"`
	ItemsID      []uint  `json:"items_id"`
	PaymentType  string  `json:"payment_type" validate:"oneof=qris cash"`
}
