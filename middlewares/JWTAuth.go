package middlewares

import (
	"net/http"
	"strings"

	"api/models"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
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
