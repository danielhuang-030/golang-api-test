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
