package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/stev029/cashier/etc/database"
	"github.com/stev029/cashier/etc/database/model"
	"github.com/stev029/cashier/etc/utils"
	httperr "github.com/stev029/cashier/http-errors"
	"github.com/stev029/cashier/models"
	"golang.org/x/crypto/bcrypt"
)

func (s *ServiceImpl) GetUsers(c *gin.Context) (*model.User, error) {
	return &model.User{}, nil
}

func (s *ServiceImpl) Login(c *gin.Context, req models.LoginRequest) (string, error) {
	var user model.User // Model User Anda

	// 1. Cari user berdasarkan username
	if row := s.db.Where("email = ?", req.Email).First(&user); row.RowsAffected == 0 && row.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password invalid"})
		return "", row.Error
	}

	// 2. Bandingkan password (Hashed vs Plain)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or password invalid"})
		return "", err
	}

	// 3. Jika OK, buat token
	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return "", err
	}

	return token, nil
}

func (s *ServiceImpl) Register(c *gin.Context, req models.RegisterRequest) (string, error) {
	var user model.User // Model User Anda

	if query := s.db.Where("email = ?", req.Email).First(&user); query.RowsAffected > 0 && query.Error == nil {
		return "", httperr.NewHttpExceptionJSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
	}

	if err := copier.Copy(&user, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return "", err
	}

	if err := s.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return "", err
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return "", err
	}

	return token, nil
}

func (s *ServiceImpl) LogoutService(ctx *gin.Context, tokenString string) error {
	// 1. Verifikasi token terlebih dahulu
	token, claims, err := utils.VerifyToken(tokenString)
	if err != nil || !token.Valid {
		return httperr.NewHttpExceptionJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	}

	// 2. Hitung sisa waktu expired (TTL)
	expTime := int64(claims["exp"].(float64))
	remainingTime := time.Unix(expTime, 0).Sub(time.Now())

	if remainingTime > 0 {
		// 3. Simpan token ke Redis dengan status 'blacklisted'
		// ctx adalah context.Background()
		err := database.RedisClient.Set(ctx, tokenString, "blacklisted", remainingTime).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
