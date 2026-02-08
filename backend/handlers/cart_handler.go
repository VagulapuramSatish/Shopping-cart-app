package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"backend/database"
	"backend/models"
)

func CreateOrAddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		ItemIDs []uint
	}
	c.BindJSON(&input)

	var cart models.Cart
	err := database.DB.Where("user_id = ?", userID).First(&cart).Error

	if err != nil {
		cart = models.Cart{
			UserID: userID,
			Status: "active",
		}
		database.DB.Create(&cart)
	}

	var items []models.Item
	database.DB.Where("id IN ?", input.ItemIDs).Find(&items)
	database.DB.Model(&cart).Association("Items").Append(items)

	c.JSON(http.StatusOK, cart)
}

func ListCarts(c *gin.Context) {
	var carts []models.Cart
	database.DB.Preload("Items").Find(&carts)
	c.JSON(http.StatusOK, carts)
}
