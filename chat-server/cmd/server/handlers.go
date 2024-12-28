package main

import (
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

/*
Handlers for all the routes present in routes.go file
- REST handling
*/

func (app *app) healthHandler(w http.ResponseWriter, request *http.Request) {
	healthResponse := map[string]string{
		"message": "Health is Ok!",
		"env":     app.appConfig.env,
		"version": version,
	}
	app.sendJSON(w, http.StatusOK, healthResponse)
}

// main chat handler which upgrade http / https connection to web socket
func (app *app) chatHandler(w http.ResponseWriter, request *http.Request) {
	webSocketConnection, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		app.errorlogger.Println("ERROR: upgrading the connection to web socket", err)
		app.internalServerErrorJSONResponse(w)
		return
	}

	app.infologger.Println("connection upgraded to Web Socket")

	client[webSocketConnection] = true

	// read the message and pass the message payload to publishChannel
	go func(connection *websocket.Conn) {
		// read the message
		for {
			var message Message
			mt, messageByte, err := connection.ReadMessage()
			if err != nil {
				app.errorlogger.Printf("ERROR: Unable to read the message from client %v: %v", webSocketConnection.RemoteAddr(), err)
				app.internalServerErrorJSONResponse(w)
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

	// create hash of password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{
			Error:   "bad_request",
			Message: err.Error(),
		})
		return
	}

	// check existing user
	user := app.userController.CheckUserExists(requestBody.Email)
	if user != nil {
		// compare password with hashedPassword
		err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(requestBody.Password))
		if err != nil {
			sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Message: err.Error(),
			})
			return
		}
		if userMap, ok := user.(map[string]interface{}); ok {
			username := userMap["username"].(string)
			email := userMap["email"].(string)
			sendJSONResponse(w, http.StatusBadRequest, SuccessResponse{
				Message: "user found",
				Data: map[string]interface{}{
					"username": username,
					"email":    email,
					"type":     "login",
				},
			})
		} else {
			sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Message: "Internal Server Error",
			})
		}
		return
	} else {
		// create new user
		user, err := app.userController.CreateNewUser(requestBody.Name, requestBody.Email, string(requestBody.Password))
		if err != nil {
			sendJSONResponse(w, http.StatusBadRequest, ErrorResponse{
				Error:   "bad_request",
				Message: err.Error(),
			})
			return
		}
		sendJSONResponse(w, http.StatusCreated, SuccessResponse{
			Message: "user created",
			Data:    user,
		})
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
