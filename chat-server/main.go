package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	Payload     string `json:"Payload"`
	PayloadType int    `json:"PayloadType"`
	// RemoteAddress string `json:"RemoteAddress"`
}

type Address struct {
	IP   string `json:"IP"`
	Port string `json:"Port"`
	Zone string `json:"Zone"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins (not secure, only for testing)
			return true
		},
	}

	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-11738.c301.ap-south-1-1.ec2.redns.redis-cloud.com:11738",
		Password: "5ZVeWbGByDGfXfxCET92Kxmz9BmboLl4", // no password set
		DB:       0,                                  // use default DB
	})

	myChannel = "common-room"

	//This map keeps track of connected WebSocket clients.
	client = make(map[*websocket.Conn]bool)

	// publish channel
	publishChannel = make(chan Message)

	// subscribe channel
	broadcastChannel = make(chan string)
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

	go publishToRedis()
	go subscribeToRedis()
	go broadcastMessages()

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
			//	message.RemoteAddress = connection.RemoteAddr().String()
			// now pass message to the broadcast channel
			publishChannel <- message
		}
	}(webSocketConnection)

}

func healthHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "my chat system's health is OK!.")
}

func broadcastMessages() {
	//	var message Message
	for {
		recieveMessage := <-broadcastChannel

		var message Message

		fmt.Println("RECIEVED IN BROADCAST:", recieveMessage)

		err := json.Unmarshal([]byte(recieveMessage), &message)
		if err != nil {
			log.Printf("ERROR: Unable to unmarshal message: %v", err)
			continue
		}

		for clientConnection := range client {
			err := clientConnection.WriteMessage(message.PayloadType, []byte(message.Payload))
			if err != nil {
				log.Println("ERROR: writing message", err)
				clientConnection.Close()
				delete(client, clientConnection)
			}
		}
	}
}

func publishToRedis() {
	for {
		payload := <-publishChannel

		payloadJson, err := json.Marshal(payload)
		if err != nil {
			log.Println("ERROR: marshalling the payload json")
		}
		err = rdb.Publish(ctx, myChannel, payloadJson).Err()
		if err != nil {
			panic(err)
		}
		log.Println("INFO: message published")
	}
}

func subscribeToRedis() {
	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, myChannel)

	// Close the subscription when we are done.
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Channel, msg.Payload)
		broadcastChannel <- msg.Payload
		log.Println("INFO: message subscribed")
	}

}
