package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

// User
type User struct {
	ID            uint `gorm:"primary_key"`
	Name          string
	Email         string
	Password      string
	ApiToken      string
	RememberToken string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Account
type Account struct {
	ID        uint `gorm:"primary_key"`
	Account   string
	Password  string
	Ip        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Setting
type Setting struct {
	ID        uint `gorm:"primary_key"`
	Skey      string
	Svalue    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// init
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

//JWTAuthMiddleware middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !validateToken(c) {
			c.Abort()
			return
		}
		c.Next()
		return
	}
}

// validate token
func validateToken(c *gin.Context) bool {
	token := c.Request.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "Unauthorized Error",
			"data":       "",
		})
		return false
	}

	// connect DB
	db, err := gorm.Open(os.Getenv("DB_CONNECTION"), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		statusCode := http.StatusBadRequest
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"message":    err.Error(),
			"data":       "",
		})
		return false
	}

	// get user token
	var user User
	db.Where("api_token = ?", token).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "Unauthorized Error",
			"data":       "",
		})
		return false
	}
	return true
}

func main() {
	router := gin.Default()
	api := router.Group("/api/v1")
	api.Use(JWTAuthMiddleware())
	{
		api.POST("/accounts", func(c *gin.Context) {
			// connect DB
			db, err := gorm.Open(os.Getenv("DB_CONNECTION"), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
			if err != nil {
				statusCode := http.StatusBadRequest
				c.JSON(statusCode, gin.H{
					"statusCode": statusCode,
					"message":    err.Error(),
					"data":       "",
				})
				return
			}

			// add account
			newAccount, err := createAccount(db, c.PostForm("account"))
			if err != nil {
				statusCode := http.StatusBadRequest
				c.JSON(statusCode, gin.H{
					"statusCode": statusCode,
					"message":    err.Error(),
					"data":       "",
				})
				return
			}
			defer db.Close()

			// Success
			statusCode := http.StatusOK
			c.JSON(statusCode, gin.H{
				"statusCode": statusCode,
				"message":    "Success",
				"data":       newAccount,
			})
		})
	}
	router.Run(":4000")
}

// get random password
func getRandomPassword(strlen int) string {
	const POOL = "abcdefghijkmnpqrstuwxyz23456789"
	password := make([]byte, strlen)
	for i := range password {
		password[i] = POOL[rand.Intn(len(POOL))]
	}
	return string(password)
}

// create account
func createAccount(db *gorm.DB, account string) (Account, error) {
	tx := db.Begin()

	// check account
	if account == "" {
		tx.Rollback()
		return Account{}, errors.New("The account is empty")
	}

	// get next IP
	var setting Setting
	tx.Find(&setting, "skey = ?", "private_ip_member")
	ip := net.ParseIP(setting.Svalue)
	ip = ip.To4()
	ip[3]++
	newIp := ip.String()
	setting.Svalue = newIp
	tx.Save(&setting)

	// add new account
	newAccount := Account{
		Account:  account,
		Password: getRandomPassword(10),
		Ip:       newIp,
	}

	if err := tx.Save(&newAccount).Error; err != nil {
		tx.Rollback()
		return Account{}, err
	}
	tx.Commit()

	return newAccount, nil
}
