package main

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

/*
JWT based authentication
*/

var tokenInvalidError = fmt.Errorf("INVALID-JWT-TOKEN")

type myJwtClaims struct {
	userId int
	jwt.RegisteredClaims
}

var secretKey = []byte("secret-key")

// craete a JWT Token
func (app *app) createJwtToken(userId int, expirationTime time.Time) (string, error) {
	claims := myJwtClaims{
		userId,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verify the JWT token
func (app *app) verifyToken(tokenString string, claims *myJwtClaims) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		app.errorlogger.Println("unable to parse jwt token")
		return false, err
	}

	if !token.Valid {
		app.infologger.Println("invalid token")
		return false, nil
	}

	return true, nil
}
