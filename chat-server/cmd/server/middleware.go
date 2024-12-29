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
