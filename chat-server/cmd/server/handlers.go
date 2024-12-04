package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

func (app *app) healthHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "my chat system's health is OK!.")
}

func (app *app) chatHandler(w http.ResponseWriter, request *http.Request) {
	webSocketConnection, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		log.Println("ERROR: upgrading the connection to web socket", err)
		return
	}

	log.Println("INFO: connection upgraded to Web Socket")

	//defer webSocketConnection.Close()

	client[webSocketConnection] = true

	// read the message and pass the message payload to publishChannel
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

func (app *app) homeHandler(w http.ResponseWriter, request *http.Request) {
	uiTemplates := "ui/index.html"
	templates, err := template.ParseFiles(uiTemplates)
	if err != nil {
		log.Println("ERROR: parsing the template files", err)
	}
	templates.Execute(w, nil)
}
