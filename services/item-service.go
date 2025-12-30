package services

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/stev029/cashier/etc/database/model"
	"github.com/stev029/cashier/models"
	"gorm.io/gorm"
)

func (s *ServiceImpl) GetItems(c *gin.Context) (*[]model.Item, error) {
	var items []model.Item

	if err := s.db.Find(&items).Error; err != nil {
		return nil, err
	}

	return &items, nil
}

func (s *ServiceImpl) GetItemByID(c *gin.Context) (*model.Item, error) {
	var item model.Item

	query := s.db.First(&item, c.Param("id"))
	if query.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	if err := query.Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *ServiceImpl) CreateItem(c *gin.Context, req models.ItemRequest) (*model.Item, error) {
	var item model.Item

	if err := copier.Copy(&item, &req); err != nil {
		return nil, err
	}

	if err := s.db.Create(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *ServiceImpl) UpdateItem(c *gin.Context, req models.ItemRequest) (*model.Item, error) {
	var item model.Item

	if err := s.db.First(&item, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	if err := copier.Copy(&item, &req); err != nil {
		return nil, err
	}

	if err := s.db.Where("id = ?", c.Param("id")).Updates(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *ServiceImpl) DeleteItem(c *gin.Context) {
	var item model.Item

	go func() {
		if err := s.db.Delete(&item, c.Param("id")).Error; err != nil {
			log.Printf("Error while deleting item: %v", err)
		}
	}()
}
