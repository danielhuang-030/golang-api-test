package main

import (
	"net/http"
	"os"

	"api/middlewares"
	"api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuth())
	api.POST("/accounts", store)

	router.Run(":" + os.Getenv("APP_PORT"))
}

func store(c *gin.Context) {
	// add account
	newAccount, err := models.CreateAccount(c.PostForm("account"))
	if err != nil {
		statusCode := http.StatusBadRequest
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"message":    err.Error(),
			"data":       "",
		})
		return
	}

	// Success
	statusCode := http.StatusOK
	c.JSON(statusCode, gin.H{
		"statusCode": statusCode,
		"message":    "Success",
		"data":       newAccount,
	})
}
