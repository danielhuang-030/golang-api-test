package main

import (
	"net/http"
	"os"
	"strings"

	"api/models"

	"github.com/gin-gonic/gin"
)

//JWTAuthMiddleware middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if !models.IsAuthByToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"message":    "Unauthorized Error",
				"data":       "",
			})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

func main() {
	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(JWTAuthMiddleware())
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
