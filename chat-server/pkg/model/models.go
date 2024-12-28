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

type Channel struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `json:"name"`
}

type Message struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	UserId    uint
	ChannelId uint
	Data      string `json:"message"`
	CreatedAt time.Time
}
