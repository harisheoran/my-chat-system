package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mainRouter := mux.NewRouter()

	mainRouter.HandleFunc("/v1/health", healthHandler)
	mainRouter.HandleFunc("/v1/chat", homeHandler)
	mainRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	log.Println("client listening on port 5000")
	err := http.ListenAndServe(":5000", mainRouter)
	if err != nil {
		log.Fatal("Unable to start the server", err)
	}
}

func healthHandler(w http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(w, "my-chat-system client health is OK!")
}

func homeHandler(w http.ResponseWriter, request *http.Request) {
	uiTemplates := "ui/main.html"
	templates, err := template.ParseFiles(uiTemplates)
	if err != nil {
		log.Println("ERROR: parsing the template files", err)
	}
	templates.Execute(w, nil)

}
