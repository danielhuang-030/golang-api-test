package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"api/middlewares"
	"api/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuth())
	api.POST("/accounts", store)
	api.DELETE("/accounts/:id", destroy)
	api.PUT("/accounts/rebuild", rebuild)

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

func destroy(c *gin.Context) {
	var err error
	if id := c.Param("id"); "" == id {
		err = errors.New("Params error")
	} else {
		id, err := strconv.Atoi(id)
		if err == nil {
			err = models.DestroyAccount(id)
		}
	}

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
		"data":       "",
	})
}

func rebuild(c *gin.Context) {
	err := models.RebuildAccountFile()
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
		"data":       "",
	})
}
