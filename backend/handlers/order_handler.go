package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"backend/database"
	"backend/models"
)

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		CartID uint
	}
	c.BindJSON(&input)

	order := models.Order{
		CartID: input.CartID,
		UserID: userID,
	}

	database.DB.Create(&order)
	c.JSON(http.StatusCreated, order)
}

func ListOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Find(&orders)
	c.JSON(http.StatusOK, orders)
}
