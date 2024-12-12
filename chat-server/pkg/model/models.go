package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	username string `json:"username"`
}

type Message struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Msg  string `json:"message"`
	Time time.Time
}
