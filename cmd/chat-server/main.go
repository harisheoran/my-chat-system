package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

/*
Role: Server
Purpose:

*/

type Message struct {
	Payload       string
	PayloadType   int
	RemoteAddress net.Addr
}

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

	//This map keeps track of connected WebSocket clients.
	client = make(map[*websocket.Conn]bool)

	// broadcast channel
	broadcast = make(chan Message)
)

func main() {
	// check redis is connected or not
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Println("ERROR: Unable to connect with Redis.", err)
	}

	// routes
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/v1/health", healthHandler)
	mainRouter.HandleFunc("/v1/chat", chatHandler)

	// write messages
	go writeMessages()

	log.Println("server starting ")
	err = http.ListenAndServe(":1316", mainRouter)
	if err != nil {
		log.Fatal("ERROR: starting the server on port 1316", err)
	}
}

func chatHandler(w http.ResponseWriter, request *http.Request) {
	webSocketConnection, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		log.Println("ERROR: upgrading the connection to web socket", err)
		return
	}

	log.Println("INFO: connection upgraded to Web Socket")

	//defer webSocketConnection.Close()

	client[webSocketConnection] = true
	fmt.Println(client)
	go func(connection *websocket.Conn) {
		// read the message
		for {
			var message Message
			mt, messageByte, err := connection.ReadMessage()
			if err != nil {
				log.Printf("ERROR: Unable to read the message from client %v: %v", webSocketConnection.RemoteAddr(), err)
				delete(client, connection)
				return
			}

			message.Payload = string(messageByte)
			message.PayloadType = mt
			message.RemoteAddress = connection.RemoteAddr()
			// now pass message to the broadcast channel
			broadcast <- message
		}
	}(webSocketConnection)

}

func healthHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "my chat system's health is OK!.")
}

func writeMessages() {
	for {
		recieveMessage := <-broadcast
		for clientConnection := range client {
			if recieveMessage.RemoteAddress == clientConnection.RemoteAddr() {
				continue
			} else {
				err := clientConnection.WriteMessage(recieveMessage.PayloadType, []byte(recieveMessage.Payload))
				if err != nil {
					log.Println("ERROR: writing message", err)
					clientConnection.Close()
					delete(client, clientConnection)
				}
			}
		}
	}
}
