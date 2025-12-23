package services

import "gorm.io/gorm"

type ServiceImpl struct {
	db *gorm.DB
}

func NewServiceImpl(db *gorm.DB) *ServiceImpl {
	return &ServiceImpl{db: db}
}
