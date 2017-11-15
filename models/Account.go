package models

import (
	"errors"
	"math/rand"
	"net"
	"time"
)

type Account struct {
	ID        uint `gorm:"primary_key"`
	Account   string
	Password  string
	Ip        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// create account
func CreateAccount(account string) (Account, error) {
	db := GetDB()
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

// get random password
func getRandomPassword(strlen int) string {
	const POOL = "abcdefghijkmnpqrstuwxyz23456789"
	password := make([]byte, strlen)
	for i := range password {
		password[i] = POOL[rand.Intn(len(POOL))]
	}
	return string(password)
}
