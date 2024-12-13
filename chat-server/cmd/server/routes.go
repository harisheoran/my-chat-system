package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
Contains all the API's routes
*/

func (app *app) router() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/v1/home", app.homeHandler)
	mainRouter.HandleFunc("/v1/health", app.healthHandler)
	mainRouter.HandleFunc("/v1/chat", app.chatHandler)
	mainRouter.HandleFunc("/v1/auth", app.authHandler)

	return mainRouter
}
