package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Username  string
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
}

type Message struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Msg  string `json:"message"`
	Time time.Time
}
