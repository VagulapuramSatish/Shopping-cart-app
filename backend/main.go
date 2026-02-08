package main

import (
	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"backend/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to DB
	database.Connect()

	// Auto migrate models
	database.DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
	)

	// Seed initial items
	SeedItems()

	// Create Gin router
	r := gin.Default()

	// === CORS Setup ===
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React dev URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// === Public routes ===
	r.POST("/users", handlers.CreateUser)
	r.POST("/users/login", handlers.Login)
	r.GET("/users", handlers.ListUsers)

	r.POST("/items", handlers.CreateItem)
	r.GET("/items", handlers.ListItems)

	// === Protected routes ===
	auth := r.Group("/")
	auth.Use(middleware.Auth()) // JWT middleware
	{
		auth.POST("/carts", handlers.CreateOrAddToCart)
		auth.GET("/carts", handlers.ListCarts)

		auth.POST("/orders", handlers.CreateOrder)
		auth.GET("/orders", handlers.ListOrders)
	}

	// Start server
	r.Run(":8080")
}

func SeedItems() {
    items := []models.Item{
        {Name: "Laptop", Status: "available"},
        {Name: "Phone", Status: "available"},
        {Name: "Headphones", Status: "available"},
        {Name: "Monitor", Status: "available"},
        {Name: "Keyboard", Status: "available"},
        {Name: "Mouse", Status: "available"},
        {Name: "Webcam", Status: "available"},
        {Name: "Printer", Status: "available"},
        {Name: "Speakers", Status: "available"},
    }

    for _, item := range items {
        if err := database.DB.Create(&item).Error; err != nil {
            println("Failed to insert item:", item.Name, err.Error())
        }
    }
    println("Seeded default items successfully!")
}
