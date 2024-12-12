package main

import (
	"fmt"
	"strings"

	"github.com/harisheoran/my-chat-system/pkg/model"
)

func (app *app) userAuth(requestBody *AuthRequestBody) (data interface{}, err error) {
	// Check for empty fields in the request body
	if requestBody.Email == "" || requestBody.Password == "" {
		fmt.Println("got null values")
		return nil, fmt.Errorf("email and/or password fields are empty")
	}

	var user model.User

	query := app.messageController.DbConnection.Where("email = ?", requestBody.Email).First(&user)

	// send existing user
	if query.RowsAffected != 0 {
		response := map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"type":     "login",
		}
		return response, nil
	}

	// create new user
	newUser := model.User{
		Name:     requestBody.Name,
		Username: strings.Split(requestBody.Email, "@")[0],
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}
	query = app.messageController.DbConnection.Create(&newUser)

	if query.Error != nil {
		fmt.Println("Got error in DB transaction:", query.Error)
		return nil, query.Error
	}
	response := map[string]interface{}{
		"username": newUser.Username,
		"email":    newUser.Email,
		"name":     newUser.Name,
		"type":     "signup",
	}
	return response, nil
}
