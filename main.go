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

var db *gorm.DB

// init
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect DB
	dbConn, err := gorm.Open(os.Getenv("DB_CONNECTION"), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatal(err.Error())
	}
	db = dbConn
}

//JWTAuthMiddleware middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if !isAuthByToken(token) {
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
	api.POST("/accounts", func(c *gin.Context) {
		// add account
		newAccount, err := createAccount(c.PostForm("account"))
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
	})

	router.Run(":" + os.Getenv("APP_PORT"))
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
func createAccount(account string) (Account, error) {
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

// is auth by token
func isAuthByToken(token string) bool {
	if token == "" {
		return false
	}

	var user User
	db.Find(&user, "api_token = ?", token)
	if user.ID == 0 {
		return false
	}
	return true
}
