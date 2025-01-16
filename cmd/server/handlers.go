package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/harisheoran/my-chat-system/internal/filter"
	"github.com/harisheoran/my-chat-system/internal/validator"
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

// main chat handler which upgrade http/https connection to web socket connection and handle chats
func (app *app) groupChatHandler(w http.ResponseWriter, request *http.Request) {

	// upgrade connection to websocket
	webSocketConnection, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "unable to upgrade the connection to web socket", err)
		return
	}

	app.infologger.Println("connection successfully upgraded to Web Socket")

	// get userID from claims
	claims, ok := request.Context().Value("userClaims").(*myJwtClaims)
	if !ok {
		app.errorlogger.Println("ERROR:", ok)
		app.internalServerErrorJSONResponse(w, "failed to retrieve the claims from request context", nil)
		return
	}

	userId := claims.UserId

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

	v := validator.New()

	// get parameters from the request
	page := app.readPaginationParameters(request.URL.Query(), "page", 1, v)
	pageSize := app.readPaginationParameters(request.URL.Query(), "page_size", 10, v)

	// check validation on page and page size
	v.Check(page > 0, "page", "must be greatar than 0")
	v.Check(pageSize > 0, "page_size", "must be greatar than 0")
	if !v.Valid() {
		app.failedValidationResponse(w, request, v.Errors)
		return
	}

	// create filter for pagination
	filter := filter.Filter{
		Page:     page,
		PageSize: pageSize,
	}

	// get messages from the database
	messagesList, err := app.messageHistory(filter)
	if err != nil {
		app.errorlogger.Println("unable to retrieve the message history ", err)
		app.internalServerErrorJSONResponse(w, "unable to retreive the message for history", err)
	}

	// send messsage list to the client
	app.sendJSON(w, http.StatusOK, messagesList)

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

// get list of all channels
func (app *app) getChannels(w http.ResponseWriter, request *http.Request) {
	channels, err := app.channelController.GetChannels()
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to retrieve channels list from the database", err)
	}

	// send channels list to the client
	err = app.sendJSON(w, http.StatusOK, channels)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "unable to send json response", err)
	}
}

/*
TODO: Create the following functionalities
- Edit channel
- Delete channl
*/
