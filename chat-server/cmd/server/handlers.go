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

			publishChannel <- message
		}
	}(webSocketConnection)
}

func (app *app) homeHandler(w http.ResponseWriter, request *http.Request) {
	uiTemplates := []string{
		"ui/index.page.tmpl",
		"ui/base.tmpl",
	}
	templates, err := template.ParseFiles(uiTemplates...)
	if err != nil {
		app.errorlogger.Println("ERROR: parsing the template files", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = templates.Execute(w, nil)
	if err != nil {
		app.errorlogger.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (app *app) authPageHandler(w http.ResponseWriter, request *http.Request) {
	mytemplate := "ui/auth.html"

	templates, err := template.ParseFiles(mytemplate)
	if err != nil {
		app.errorlogger.Println("ERROR: parsing the template files", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = templates.Execute(w, nil)
	if err != nil {
		app.errorlogger.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (app *app) authHandler(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		app.errorlogger.Print("unable to parse the form", err)

	}
}
