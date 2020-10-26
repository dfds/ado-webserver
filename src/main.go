package main

import (
	"os"
	"net/http"
	"log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	adoHandlers "ado-pipeline/handlers"
)

func main() {
	// Routing
	r := mux.NewRouter()
	r.HandleFunc("/builds", adoHandlers.GetBuildsHandler)

	r.Use(adoHandlers.CorsHandler(r))

	// Launch HTTP server
	println("HTTP server listening on :8080")
	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))); err != nil {
		log.Fatal(err)
	}
}