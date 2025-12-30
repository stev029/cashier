package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/stev029/cashier/middlewares"
	"gorm.io/gorm"
)

func Controller(routes *gin.Engine, db *gorm.DB) {
	// Items Routes
	itemRoute := routes.Group("/items")
	itemRoute.Use(middlewares.AuthMiddleware())
	itemRoute.Use(middlewares.GroupPermission("staff"))
	{
		itemController := NewItemController(db)
		itemRoute.POST("/", itemController.CreateItem)
		itemRoute.GET("/", itemController.GetItems)
		itemRoute.GET("/:id", itemController.GetItemByID)
		itemRoute.PUT("/:id", itemController.UpdateItem)
		itemRoute.DELETE("/:id", itemController.DeleteItem)
	}

	// Auth Routes
	authRoute := routes.Group("/auth")
	{
		authController := NewAuthController(db)
		authRoute.POST("/register", authController.Register)
		authRoute.POST("/login", authController.Login)
		authRoute.POST("/logout", middlewares.AuthMiddleware(), authController.Logout)
	}

	// User Routes
	userRoute := routes.Group("/users")
	{
		userController := NewUserController(db)
		userRoute.GET("/me", middlewares.AuthMiddleware(), userController.Me)
		userRoute.PUT("/:id/update", middlewares.AuthMiddleware(), middlewares.GroupPermission("superuser"), userController.UpdateUser)
		userRoute.GET("/:id", userController.GetUserByID)
	}

	// Transaction Routes
	transactionRoute := routes.Group("/transactions")
	{
		transactionController := NewTransactionController(db)
		transactionRoute.POST("/", middlewares.AuthMiddleware(), middlewares.GroupPermission("staff"), transactionController.CreateTransaction)
		transactionRoute.POST("/webhook", transactionController.WebhookTransaction)
	}
}
