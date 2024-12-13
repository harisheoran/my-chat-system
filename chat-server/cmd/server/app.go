package main

import (
	"log"

	postgre "github.com/harisheoran/my-chat-system/pkg/model/postgre"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

/*
contains all the classes used in the API
*/

// main application class, used to share the dependencies across whole application
type app struct {
	infologger        *log.Logger
	errorlogger       *log.Logger
	redisConnection   *redis.Client
	messageController postgre.MessageController
	userController    postgre.UserController
	kafkaProducer     *kafka.Writer
}

type Message struct {
	Payload     string `json:"Payload"`
	PayloadType int    `json:"PayloadType"`
	// RemoteAddress string `json:"RemoteAddress"`
}
