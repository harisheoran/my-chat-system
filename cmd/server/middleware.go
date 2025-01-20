package main

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

/*
Contains all the middlewares
*/

// authentication middleware
func (app *app) CheckAutheticationMiddleware(handler http.Handler) http.Handler {
	nextHandler := func(w http.ResponseWriter, request *http.Request) {
		// get the cookie
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

		// get JWT string from the cookie
		jwtToken := cookie.Value

		myClaims := &myJwtClaims{}

		isTokenVerified, myClaims, err := app.verifyToken(jwtToken, myClaims)
		if err != nil {
			if err == jwt.ErrTokenSignatureInvalid {
				app.serverErrorJsonResponse(w, http.StatusUnauthorized, ErrorResponse{
					Message: "Not authorized, Invalid Token",
				})
			}
			app.serverErrorJsonResponse(w, http.StatusBadRequest, ErrorResponse{
				Message: "Bad Request",
			})
			return
		}

		if !isTokenVerified {
			app.sendJSON(w, http.StatusUnauthorized, NotFoundResponse{
				Message: "Session expired",
			})
			return
		}

		ctx := request.Context()
		ctx = context.WithValue(request.Context(), "userClaims", myClaims)

		handler.ServeHTTP(w, request.WithContext(ctx))
	}

	return http.HandlerFunc(nextHandler)
}

// CORS middleware
func (app *app) corsMiddleware(nextHandler http.Handler) http.Handler {

	handler := func(w http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		allowedOrigin := "http://localhost:5173" // Change to your frontend URL

		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
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
