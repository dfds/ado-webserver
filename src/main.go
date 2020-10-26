package main

import (
	"os"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	adoHandlers "ado-pipeline/handlers"
	identity "ado-pipeline/identity"
)

func main() {	
	// Acquire ADO token from environment variable "ADO_TOKEN"
	//token := os.Getenv("ADO_TOKEN")

	//if token == "" {
		//panic("A valid Azure DevOps access token needs to be set in the environment variable: 'ADO_TOKEN'")
	//}

	identity.AcquireTokenClientSecret();

	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/builds", adoHandlers.GetBuildsHandler)

	r.Use(adoHandlers.CorsHandler(r))
	//r.Use(adoHandlers.OAuth2Handler(r))

	// Launch HTTP server
	println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))); err != nil {
		log.Fatal(err)
	}
}