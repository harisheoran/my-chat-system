package main

import (
	"encoding/json"
	"net/http"
)

/*
Helper functions to handle common and repeated tasks
*/
func sendJSONResponse(w http.ResponseWriter, statusCode int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// send json responsei
func (app *app) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)

	return nil
}

func (app *app) readJSON(request *http.Request, target interface{}) error {

	err := json.NewDecoder(request.Body).Decode(&target)
	if err != nil {
		return err
	}

	return nil
}

/*
Error JSON response helpers
*/

// server error response in JSON
func (app *app) serverErrorJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	app.errorlogger.Println(data)
	err := app.sendJSON(w, statusCode, "")
	if err != nil {
		app.errorlogger.Println("Unable to send internal server error response ", err)
		w.WriteHeader(500)
	}
}

// internal server error response in JSON
func (app *app) internalServerErrorJSONResponse(w http.ResponseWriter) {
	message := "The server encountered an Internal Error and could not process the request."
	app.serverErrorJsonResponse(w, http.StatusInternalServerError, message)
}
