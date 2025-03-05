package main

import (
	"assignment1/handlers"
	"assignment1/utils"
	"log"
	"net/http"
	"os"
)

/*
 * Entry point for the service. Initializes the server, registers handlers,
 * and starts listening for incoming requests.
 */
func main() {
	// Sets the port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		port = utils.DEFAULT_PORT

	}
	// Handler endpoints
	http.HandleFunc(utils.DEFAULT_PATH, handlers.DefaultHandler)		// Root endpoint
	http.HandleFunc(utils.INFO_PATH, handlers.InfoHandler)				// Country info encpoint
	http.HandleFunc(utils.POPULATION_PATH, handlers.PopulationHandler)	// Population data endpoint
	http.HandleFunc(utils.STATUS_PATH, handlers.StatusHandler)			// Status endpoint

	// Start server
	log.Println("Starting server on port " + port + " .... ")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
