package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

/*
Handlers for all the routes present in routes.go file
- REST handling
*/

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

func (app *app) authHandler(w http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST" {
		sendJSONResponse(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error:   "method_not_allowed",
			Message: "method not allowed",
		})
		return
	}

	reqBodyContent, err := io.ReadAll(request.Body)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
			Error:   err.Error(),
			Message: "Internal Server Error",
		})
		return
	}
	defer request.Body.Close()
	var requestBody AuthRequestBody
	err = json.Unmarshal(reqBodyContent, &requestBody)
	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
			Error:   err.Error(),
			Message: "Internal Server Error",
		})
		return
	}
	data, err := app.userAuth(&requestBody)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})
		return
	}
	sendJSONResponse(w, http.StatusOK, SuccessResponse{
		Data:    data,
		Message: "authentication success",
	})
}
