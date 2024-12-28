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

	mainRouter.HandleFunc("/health", app.healthHandler)
	mainRouter.HandleFunc("/signup", app.signupHandler).Methods("POST")
	mainRouter.HandleFunc("/login", app.loginHandler).Methods("POST")

	appRouter := mainRouter.PathPrefix("/v1").Subrouter()
	appRouter.HandleFunc("/history", app.messageHistoryHandler).Methods("GET")
	appRouter.HandleFunc("/home", app.homeHandler)
	appRouter.HandleFunc("/chat", app.chatHandler)
	appRouter.HandleFunc("/logout", app.logoutHandler)
	appRouter.Use(app.CheckAutheticationMiddleware)
	mainRouter.HandleFunc("/v1/online-users/add/{userId}", app.addOnlineUser).Methods("POST")
	mainRouter.HandleFunc("/v1/online-users/remove/{userId}", app.removeOnlineUser).Methods("POST")
	mainRouter.HandleFunc("/v1/online-users/get-count", app.getOnlineUsersCount).Methods("GET")

	return mainRouter
}
