package main

import (
	"log"

	"github.com/redis/go-redis/v9"
)

type app struct {
	infologger      *log.Logger
	errorlogger     *log.Logger
	redisConnection *redis.Client
}

type Message struct {
	Payload     string `json:"Payload"`
	PayloadType int    `json:"PayloadType"`
	// RemoteAddress string `json:"RemoteAddress"`
}
