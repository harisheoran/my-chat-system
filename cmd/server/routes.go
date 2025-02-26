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
	appRouter.HandleFunc("/channel/{channelid}", app.groupChatHandler)
	appRouter.HandleFunc("/logout", app.logoutHandler).Methods("GET")
	appRouter.HandleFunc("/create-channel", app.createChannelHandler).Methods("POST")
	appRouter.HandleFunc("/getchannels", app.getChannels).Methods("GET")
	appRouter.HandleFunc("/online-users/add/{userId}", app.addOnlineUser).Methods("POST")
	appRouter.HandleFunc("/online-users/remove/{userId}", app.removeOnlineUser).Methods("POST")
	appRouter.HandleFunc("/online-users/get-count", app.getOnlineUsersCount).Methods("GET")
	appRouter.Use(app.CheckAutheticationMiddleware)

	return app.corsMiddleware(mainRouter)
}
