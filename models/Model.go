package models

import (
  "fmt"
  "log"
  "os"

  _ "github.com/go-sql-driver/mysql"
  "github.com/jinzhu/gorm"
  "github.com/joho/godotenv"
)

// db
var db *gorm.DB

// init
func init() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

  // connect db
  dbConn, err := gorm.Open(os.Getenv("DB_CONNECTION"), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
  if err != nil {
    log.Fatal(err.Error())
  }
  db = dbConn
}

// get db
func GetDB() *gorm.DB {
  return db
}
