package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/harisheoran/my-chat-system/pkg/model"
)

// user sign up
func (app *app) signupHandler(w http.ResponseWriter, request *http.Request) {
	user := model.User{}
	// read the json from request body
	err := app.readJSON(request, &user)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "unable to read the request payload for user signup ", err)
	}

	// check user exist or not
	userExist, err := app.userController.CheckUserExist(user.Email)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to query the database", err)
		return
	}

	if userExist {
		response := NotFoundResponse{
			"User already exist, please login",
		}
		app.sendJSON(w, http.StatusOK, response)
		return
	}

	// save the user info
	err = app.userController.CreateNewUser(user)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to save the user credentials to the database", err)
	}

	// send success response after signup
	successResponse := SuccessResponse{
		Message: "User signed up successfully",
	}
	err = app.sendJSON(w, http.StatusOK, successResponse)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to send the success response of saving user signup", err)
	}
}

// user login handler
func (app *app) loginHandler(w http.ResponseWriter, request *http.Request) {
	// decode the request
	loginPayload := LoginRequestPayload{}
	err := app.readJSON(request, &loginPayload)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "unable to decode the login payload", err)
	}

	// check the user exists or not
	userExist, err := app.userController.CheckUserExist(loginPayload.Email)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to check the existing user", err)
	}
	if !userExist {
		err = app.sendJSON(w, http.StatusNotFound, NotFoundResponse{
			Message: "User does not exist",
		})
		if err != nil {
			app.internalServerErrorJSONResponse(w, "failed to send json response", err)
		}
		return
	}

	// check the authenticty of payload credentials
	userId, err := app.userController.Authenticate(loginPayload.Email, loginPayload.Password)
	if err != nil || userId < 0 {
		err := app.sendJSON(w, http.StatusNotFound, NotFoundResponse{
			Message: "User Email or Password does not match.",
		})
		if err != nil {
			app.internalServerErrorJSONResponse(w, "failed to send json response", err)
		}
	}

	// create JWT token and send to the client
	expirationTime := time.Now().Add(5 * time.Minute)
	token, err := app.createJwtToken(userId, expirationTime)
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to sign the JWT Token", err)
	}

	// set a jwt token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: expirationTime,
		Secure:  false,
		Path:    "/",
	})

	userIdString := strconv.Itoa(userId)

	// set a cookie for saving userID
	http.SetCookie(w, &http.Cookie{
		Name:    "userId",
		Value:   userIdString,
		Expires: expirationTime,
		Secure:  false,
		Path:    "/",
	})

	// send successfull login response
	err = app.sendJSON(w, http.StatusOK, SuccessResponse{
		Message: "Login succeeded",
	})
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to send json response", err)
	}
}

// logout handler
func (app *app) logoutHandler(w http.ResponseWriter, request *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now(),
		Secure:  false,
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "userId",
		Value:   "",
		Expires: time.Now(),
		Secure:  false,
		Path:    "/",
	})
	err := app.sendJSON(w, http.StatusOK, SuccessResponse{
		Message: "Logout successfully",
	})
	if err != nil {
		app.internalServerErrorJSONResponse(w, "failed to send json response", err)
	}
}
