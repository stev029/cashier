package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	UserID      uint           `json:"user_id"`
	Amount      int            `json:"amount"`
	OrderID     string         `json:"order_id"`
	Status      string         `json:"status"`
	PaymentType string         `json:"payment_type"`
	PaymentCode string         `json:"payment_code"`
	CreatedAt   int64          `json:"created_at"`
	UpdatedAt   int64          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	User User   `gorm:"foreignKey:UserID"`
	Item []Item `gorm:"many2many:transaction_items"`
}
