package services

import "gorm.io/gorm"

type ServiceImpl struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *ServiceImpl {
	return &ServiceImpl{db: db}
}
