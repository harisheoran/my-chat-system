package main

import (
	"log"
	"time"

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
	messageController *postgre.MessageController
	userController    *postgre.UserController
	channelController *postgre.ChannelController
	kafkaProducer     *kafka.Writer
	kafkaConsumer     *kafka.Reader
	appConfig         *AppConfig
	kafkaUrl          string
	jwtSecretKey      string
}

type AppConfig struct {
	port int
	env  string
}

type Message struct {
	PayloadType int
	UserId      uint
	ChannelId   uint
	Data        string `json:"message"`
	CreatedAt   time.Time
}

type LoginRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
