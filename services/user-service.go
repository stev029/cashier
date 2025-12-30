package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/stev029/cashier/etc/database/model"
	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
)

func (s *ServiceImpl) GetUserByID(ctx *gin.Context) (*model.User, error) {
	var user model.User

	if err := s.db.Omit("password", "created_at", "updated_at", "deleted_at").First(&user, ctx.Param("id")).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *ServiceImpl) GetUser(ctx *gin.Context) (*model.User, error) {
	var user model.User
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, httperr.Unauthorized
	}

	if err := s.db.Omit("password", "created_at", "updated_at", "deleted_at").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *ServiceImpl) UpdateUser(ctx *gin.Context, req *models.UserRequest) error {
	var user model.User
	var group []model.Group
	userID := ctx.Param("id")

	if err := s.db.Preload("Groups").
		Omit("created_at", "updated_at", "deleted_at").
		First(&user, userID).Error; err != nil {
		return err
	}

	if err := s.db.Find(&group, req.GroupID).Error; err != nil {
		return err
	}

	user.Groups = group

	if err := copier.Copy(&user, &req); err != nil {
		return err
	}

	// log.Println("DB: %#v", s.db.Session(&gorm.Session{FullSaveAssociations: true}))

	if err := s.db.Updates(&user).Error; err != nil {
		return err
	}

	return nil
}
