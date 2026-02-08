package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"backend/database"
	"backend/handlers"
	"backend/models"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Shopping Cart API", func() {
	var router *gin.Engine

	BeforeSuite(func() {
		database.Connect()
		database.DB.Exec("DELETE FROM users")
		database.DB.Exec("DELETE FROM items")
		database.DB.Exec("DELETE FROM carts")
		database.DB.Exec("DELETE FROM orders")
		router = gin.Default()

		// User routes
		router.POST("/users", handlers.CreateUser)
		router.POST("/users/login", handlers.Login)

		// Item routes
		router.POST("/items", handlers.CreateItem)
		router.GET("/items", handlers.ListItems)

		// Cart routes
		auth := router.Group("/")
		auth.Use(handlers.TestAuthMiddleware()) // optional, see note
		{
			auth.POST("/carts", handlers.CreateOrAddToCart)
			auth.GET("/carts", handlers.ListCarts)
			auth.POST("/orders", handlers.CreateOrder)
			auth.GET("/orders", handlers.ListOrders)
		}
	})

	Describe("User signup & login", func() {
		It("should signup a new user", func() {
			user := map[string]string{"username": "alice", "password": "secret"}
			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusCreated))

			var u models.User
			err := json.Unmarshal(resp.Body.Bytes(), &u)
			Expect(err).To(BeNil())
			Expect(u.Username).To(Equal("alice"))
		})

		It("should login the user and return a token", func() {
			user := map[string]string{"username": "alice", "password": "secret"}
			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusOK))

			var res map[string]string
			err := json.Unmarshal(resp.Body.Bytes(), &res)
			Expect(err).To(BeNil())
			Expect(res).To(HaveKey("token"))
		})
	})

	Describe("Cart & Items", func() {
		var token string
		var itemID uint

		BeforeEach(func() {
			// Login to get token
			user := map[string]string{"username": "alice", "password": "secret"}
			body, _ := json.Marshal(user)
			req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			var res map[string]string
			json.Unmarshal(resp.Body.Bytes(), &res)
			token = res["token"]
		})

		It("should create an item", func() {
			item := map[string]string{"name": "Book", "status": "available"}
			body, _ := json.Marshal(item)
			req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusCreated))

			var it models.Item
			json.Unmarshal(resp.Body.Bytes(), &it)
			Expect(it.Name).To(Equal("Book"))
			itemID = it.ID
		})

		It("should create a cart and add items", func() {
			cart := map[string]interface{}{
				"name":     "Alice Cart",
				"item_ids": []uint{itemID},
			}
			body, _ := json.Marshal(cart)
			req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusCreated))
		})

		It("should convert cart into an order", func() {
			cart := models.Cart{}
			database.DB.First(&cart)
			payload := map[string]uint{"cart_id": cart.ID}
			body, _ := json.Marshal(payload)

			req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			Expect(resp.Code).To(Equal(http.StatusCreated))
		})
	})
})
