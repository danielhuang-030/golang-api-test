package models

import (
	"time"
)

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

// is auth by token
func IsAuthByToken(token string) bool {
	if token == "" {
		return false
	}

	var user User
	GetDB().Find(&user, "api_token = ?", token)
	if user.ID == 0 {
		return false
	}
	return true
}
