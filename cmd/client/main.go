package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-11738.c301.ap-south-1-1.ec2.redns.redis-cloud.com:11738",
		Password: "5ZVeWbGByDGfXfxCET92Kxmz9BmboLl4", // no password set
		DB:       0,                                  // use default DB
	})

	myChannel = "chat"

	serverAddress = "http://127.0.0.1:1313/chat"
)

func main() {

	pubsub := rdb.Subscribe(ctx, myChannel)
	defer pubsub.Close()

	mainRouter := mux.NewRouter()

	mainRouter.HandleFunc("/connectserver", connectServerHandler)

	http.ListenAndServe(":1313", mainRouter)

}

func connectServerHandler(w http.ResponseWriter, request *http.Request) {

	response, err := http.Get(serverAddress)
	if err != nil {
		log.Println("ERROR: connecting to server", err)
	}

	body := response.Body

	fmt.Println("RESPONSE", body)

	fmt.Fprintf(w, "HELLO", body)
}
