package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"backend/database"
	"backend/models"
)

func CreateItem(c *gin.Context) {
	var item models.Item
	c.BindJSON(&item)

	database.DB.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func ListItems(c *gin.Context) {
	var items []models.Item
	database.DB.Find(&items)
	c.JSON(http.StatusOK, items)
}
