package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

/*
Contains all the middlewares
*/

func (app *app) CheckAutheticationMiddleware(handler http.Handler) http.Handler {
	nextHandler := func(w http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				app.sendJSON(w, http.StatusUnauthorized, NotAuthorizedResponse{
					Message: "Not Authorized",
				})
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		jwtToken := cookie.Value

		myClaims := &myJwtClaims{}

		isTokenVerified, err := app.verifyToken(jwtToken, myClaims)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !isTokenVerified {
			app.sendJSON(w, http.StatusUnauthorized, NotFoundResponse{
				Message: "Session expired",
			})
			return
		}

		handler.ServeHTTP(w, request)
	}

	return http.HandlerFunc(nextHandler)
}

// CORS middleware
func (app *app) corsMiddleware(nextHandler http.Handler) http.Handler {

	handler := func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if request.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		nextHandler.ServeHTTP(w, request)
	}

	return http.HandlerFunc(handler)
}
