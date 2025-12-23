package database

import (
	"log"
	"os"
	"time"

	"github.com/stev029/cashier/etc/database/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DBConnect() *gorm.DB {
	if DB != nil {
		return DB
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Error while connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error while accessing DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Minute * 10)

	DB = db

	return db
}

func InitModel() error {
	if DB == nil {
		return gorm.ErrInvalidDB
	}

	err := DB.AutoMigrate(
		&model.Item{},
		&model.Category{},
		&model.User{},
		&model.Group{},
		&model.Permission{},
	)

	if err != nil {
		return err
	}

	return nil
}
