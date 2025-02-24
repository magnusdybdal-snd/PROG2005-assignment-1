package main

import (
	"assignment1/handlers"
	"assignment1/utils"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = "8080"

	}
	// Handler endpoints
	http.HandleFunc(utils.DEFAULT_PATH, handlers.DefaultHandler)
	http.HandleFunc(utils.INFO_PATH, handlers.InfoHandler)
	http.HandleFunc(utils.POPULATION_PATH, handlers.PopulationHandler)
	http.HandleFunc(utils.STATUS_PATH, handlers.StatusHandler)

	// Start server
	log.Println("Starting server on port " + port + " .... ")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
