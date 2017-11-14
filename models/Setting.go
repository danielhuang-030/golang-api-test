package models

import (
	"time"
)

type Setting struct {
	ID        uint `gorm:"primary_key"`
	Skey      string
	Svalue    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
