package main

import (
	"log"

	postgre "github.com/harisheoran/my-chat-system/pkg/model/postgre"
	"github.com/redis/go-redis/v9"
)

type app struct {
	infologger        *log.Logger
	errorlogger       *log.Logger
	redisConnection   *redis.Client
	messageController postgre.MessageController
}

type Message struct {
	Payload     string `json:"Payload"`
	PayloadType int    `json:"PayloadType"`
	// RemoteAddress string `json:"RemoteAddress"`
}
