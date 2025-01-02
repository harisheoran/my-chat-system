package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/harisheoran/my-chat-system/pkg/model"
)

/*
Handlers for all the routes of the API
*/

// health check
func (app *app) healthHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Println("HERE")
	healthResponse := map[string]string{
		"message": "Health is Ok!",
		"env":     app.appConfig.env,
		"version": version,
	}
	app.sendJSON(w, http.StatusOK, healthResponse)
}

// main chat handler which upgrade http/https connection to web socket connection
func (app *app) groupChatHandler(w http.ResponseWriter, request *http.Request) {

	// upgrade connection to websocket
	webSocketConnection, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		app.internalServerErrorJSONResponse(w, " unable to upgrade the connection to web socket", err)
		return
	}

	app.infologger.Println("connection successfully upgraded to Web Socket")

	// get userId from cookies
	userId, err := app.getUserIdFromCookie(request)
	if err == cookieNotFoundError {
		app.errorlogger.Println("user id cookie not found", err)
		return
	} else if err != nil {
		app.infologger.Println("failed to get userId cookie", err)
		return
	}

	// get channel id from path
	vars := mux.Vars(request)
	channelIdVar := vars["channelid"]
	channelId, err := strconv.ParseInt(channelIdVar, 10, 64)

	client[webSocketConnection] = true

	// read the message and pass the message payload to publishChannel
	go func(connection *websocket.Conn) {
		// read the message
		for {
			var message Message
			messageType, messageByte, err := connection.ReadMessage()
			if err != nil {
				app.internalServerErrorJSONResponse(w, "unable to read the message from web socker client", err)
				delete(client, connection)
				return
			}

			message.PayloadType = messageType
			message.UserId = uint(userId)
			message.Data = string(messageByte)
			message.ChannelId = uint(channelId)
			message.CreatedAt = time.Now()

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

/*
message history handler
*/
func (app *app) messageHistoryHandler(w http.ResponseWriter, request *http.Request) {
	err := app.messageHistory()

	if err != nil {
		app.errorlogger.Println("unable to retrieve the message history ", err)
		// send internal server error response
	}
}

func (app *app) addOnlineUser(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	userID := vars["userId"]

	go app.addToOnlineUsers(userID)

	sendJSONResponse(w, http.StatusOK, SuccessResponse{
		Message: "user added to online list",
		Data:    nil,
	})
}

func (app *app) removeOnlineUser(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	userID := vars["userId"]

	go app.removeFromOnlineUsers(userID)

	sendJSONResponse(w, http.StatusOK, SuccessResponse{
		Message: "user removed from online list",
		Data:    nil,
	})
}

func (app *app) getOnlineUsersCount(w http.ResponseWriter, request *http.Request) {
	count, err := app.countOnlineUsers()

	if err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_server_error",
			Message: "Internal Server Error",
		})
		return
	}

	sendJSONResponse(w, http.StatusOK, SuccessResponse{
		Message: "user added to online list",
		Data:    count,
	})
}

/*
create channel handler
*/
func (app *app) createChannelHandler(w http.ResponseWriter, request *http.Request) {
	// read the json payloaf
	channelPayload := model.Channel{}
	err := app.readJSON(request, &channelPayload)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to read the json payload", err)
	}

	// insert into the database
	err = app.channelController.InsertChannel(&channelPayload)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to save the channel data into the database", err)
	}

	// send success response
	err = app.sendJSON(w, http.StatusOK, SuccessResponse{
		Message: "Channel created successfully",
	})
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to send the json response", err)
	}
}

/*
TODO: Create the following functionalities
- Edit channel
- Delete channl
- View all the channels
*/
