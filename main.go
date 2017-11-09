package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

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

// // Response
// type Response struct {
// 	statusCode int
// 	message    string
// }
//
// func getResponse() Response {
// 	response := Response{}
// 	response.statusCode = http.StatusBadRequest
// 	response.message = "Error"
// 	return response
// }

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	router := gin.Default()
	router.POST("/api/v1/accounts", func(c *gin.Context) {
		// connect DB
		db, err := gorm.Open(os.Getenv("DB_CONNECTION"), fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE")))
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

		// var account Account
		// db.Find(&account, "id = ?", 3)
		// fmt.Printf("%+v", whatever)
		// account.Account = "changeByGo"
		// db.Save(&account)

		statusCode := http.StatusOK
		c.JSON(statusCode, gin.H{
			"statusCode": statusCode,
			"message":    "Success",
			"data":       newAccount,
		})
	})
	router.Run(":4000")
}

// func checkErr(err error, response Response) Response {
// 	if err != nil {
// 		response.message = err.Error()
// 		// fmt.Printf("%v", err)
// 		// panic(err)
// 	}
// 	return response
// }
// func render

func getRandomPassword(strlen int) string {
	const POOL = "abcdefghijkmnpqrstuwxyz23456789"
	password := make([]byte, strlen)
	for i := range password {
		password[i] = POOL[rand.Intn(len(POOL))]
	}
	return string(password)
}

// func getNextIp(db *gorm.DB, skey string) string {
// 	if skey == "" {
// 		skey = "private_ip_member"
// 	}
// 	var setting Setting
// 	db.Find(&setting, "skey = ?", skey)
// 	ip := net.ParseIP(setting.Svalue)
// 	ip = ip.To4()
// 	ip[3]++
//
// 	newIp := ip.String()
// 	setting.Svalue = newIp
// 	db.Save(&setting)
//
// 	return newIp
// }

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

/*
package main

import (
  "database/sql"
  "fmt"

  _ "github.com/go-sql-driver/mysql"
)

func main() {
  db, err := sql.Open("mysql", "root:@/vpn?charset=utf8mb4")
  checkErr(err)
  rows, err := db.Query("SELECT id, account FROM accounts WHERE id = 1")
  checkErr(err)
  fmt.Println(rows)

  for rows.Next() {
    var uid int
    var account string

    err = rows.Scan(&uid, &account)
    checkErr(err)

    fmt.Println(uid)
    fmt.Println(account)
  }
}
*/
