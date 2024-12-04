package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

/*
Role: Server
Purpose:
*/

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		/*CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (not secure, only for testing)
		return true
		},*/
	}

	myChannel = "common-room"

	ctx = context.Background()

	//This map keeps track of connected WebSocket clients.
	client = make(map[*websocket.Conn]bool)

	// publish channel
	publishChannel = make(chan Message)

	// subscribe channel
	broadcastChannel = make(chan string)
)

func main() {
	// create two loggers for info and error
	errorlogger := log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime)
	infologger := log.New(os.Stdin, "INFO", log.Ldate|log.Ltime|log.Lshortfile)

	// load .env file
	err := godotenv.Load()
	if err != nil {
		errorlogger.Println("can't read the env files")
	}

	// establish the redis connection
	redis_db, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 2, 64)
	if err != nil {
		errorlogger.Println("unable to parse the redis db value from .env file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       int(redis_db),
	})
	// check redis is connected or not
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		errorlogger.Println("unable to connect to redis", err)
	}

	app := app{
		infologger:      infologger,
		errorlogger:     errorlogger,
		redisConnection: rdb,
	}

	go app.publishToRedis()
	go app.subscribeToRedis()
	go app.broadcastMessages()

	// start the server
	port := fmt.Sprintf(":%s", "1316")
	server := &http.Server{
		Addr:     port,
		ErrorLog: app.errorlogger,
		Handler:  app.router(),
	}

	app.infologger.Println("Server is starting on the port ", port)
	err = server.ListenAndServe()
	if err != nil {
		app.errorlogger.Fatal("unable to start the server on port 1316", err)
	}
}
