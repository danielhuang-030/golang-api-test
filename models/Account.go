package models

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/jinzhu/gorm"
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
	// tx := GetTx()
	tx := GetDB().Begin()

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

	// append account info
	file := os.Getenv("VPN_ACCOUNT_FILE")
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return newAccount, err
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("\"%s\" l2tpd \"%s\" %s\n", newAccount.Account, newAccount.Password, newAccount.Ip)); err != nil {
		tx.Rollback()
		return Account{}, err
	}
	tx.Commit()

	return newAccount, nil
}

// destroy account
func DestroyAccount(id int) error {
	// tx := GetTx()
	tx := GetDB().Begin()

	// check id
	if id == 0 {
		return errors.New("Params error")
	}

	// get account
	var account Account
	tx.First(&account, id)
	if account.ID == 0 {
		tx.Rollback()
		return errors.New("This account does not exist")
	}

	// destroy account
	if err := tx.Delete(&account).Error; err != nil {
		tx.Rollback()
		return err
	}

	// rebuild account file
	if err := RebuildAccountFile(tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

// get random password
func getRandomPassword(strlen int) string {
	rand.Seed(time.Now().UnixNano())
	const POOL = "abcdefghijkmnpqrstuwxyz23456789"
	password := make([]byte, strlen)
	for i := range password {
		password[i] = POOL[rand.Intn(len(POOL))]
	}
	return string(password)
}

// rebuild account file
func RebuildAccountFile(tx ...*gorm.DB) error {
	file := os.Getenv("VPN_ACCOUNT_FILE")
	now := time.Now()
	datetime := now.Format("20060102150405")

	// backup
	backupPath := fmt.Sprintf("/%s/%s/%s/", os.Getenv("VPN_ACCOUNT_BACKUP_PATH"), now.Format("2006"), now.Format("01"))
	os.MkdirAll(backupPath, os.ModePerm)
	fileInfo, err := os.Stat(file)
	if err != nil {
		return err
	}
	err = copyFile(file, fmt.Sprintf("%s/%s_%s", backupPath, fileInfo.Name(), datetime))
	if err != nil {
		return err
	}

	// get all accounts
	var accounts []Account
	if len(tx) > 0 {
		tx[0].Find(&accounts)
	} else {
		GetDB().Find(&accounts)
	}
	if 0 == len(accounts) {
		return errors.New("The account list is empty")
	}

	// create file
	fileBuilding := fmt.Sprintf("%s_building_%s", file, datetime)
	f, err := os.Create(fileBuilding)
	if err != nil {
		return err
	}
	defer f.Close()

	// append account info
	for i := 0; i < len(accounts); i++ {
		f, err := os.OpenFile(fileBuilding, os.O_APPEND|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		// fmt.Printf("%v, %T", accounts[i], accounts[i])
		if _, err = f.WriteString(fmt.Sprintf("\"%s\" l2tpd \"%s\" %s\n", accounts[i].Account, accounts[i].Password, accounts[i].Ip)); err != nil {
			return err
		}
	}

	// delte origin file
	err = os.Remove(file)
	if err != nil {
		return err
	}

	// rename file
	err = os.Rename(fileBuilding, file)
	if err != nil {
		return err
	}
	return nil
}

// copy file
func copyFile(from string, to string) error {
	srcFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
