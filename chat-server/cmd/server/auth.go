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
	UserId int
	jwt.RegisteredClaims
}

// craete a JWT Token
func (app *app) createJwtToken(userId int, expirationTime time.Time) (string, error) {
	var secretKey = []byte(app.jwtSecretKey)
	claims := myJwtClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// create token with claims and signing algorithim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// create the JWT string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verify the JWT token
func (app *app) verifyToken(tokenString string, claims *myJwtClaims) (bool, *myJwtClaims, error) {
	var secretKey = []byte(app.jwtSecretKey)

	// parse the JWT string and store the results in claims
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		app.errorlogger.Println("unable to parse jwt token")
		return false, claims, err
	}

	if !token.Valid {
		app.infologger.Println("invalid token")
		return false, claims, nil
	}

	return true, claims, nil
}
