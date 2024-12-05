package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *app) router() http.Handler {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/v1/home", app.homeHandler)
	mainRouter.HandleFunc("/v1/health", app.healthHandler)
	mainRouter.HandleFunc("/v1/chat", app.chatHandler)

	return mainRouter
}
