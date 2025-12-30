package model

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name        string   `gorm:"type:varchar(100);not null"`
	Description string   `gorm:"type:text"`
	Quantity    int      `gorm:"type:int;not null"`
	Price       float64  `gorm:"type:decimal(10,2);not null"`
	Category_id uint     `gorm:"column:category_id"`
	Category    Category `gorm:"foreignKey:Category_id"`
}

type Category struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"`
	Description string `gorm:"type:text"`
}
