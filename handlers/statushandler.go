package handlers

import (
	"assignment1/utils" // const.go, func.go, struct.go
	"encoding/json"
	"net/http"
	"time"
)

// Time since service has started
var startTime = time.Now()

// URL for APIs. the exact endpoint is not important
var countriesNowURL = "http://129.241.150.113:3500/api/v0.1/countries/iso"
var RESTCountriesURL = "http://129.241.150.113:8080/v3.1/alpha/no"

/*
 * Handler for requests to the status entry point
 *
 * Checks the status of the CountriesNow and REST Countries APIs
 * Calculates the uptime of the service and gives a JSON response.
 *
 * Parameters:
 *	- w: The http.ResponseWriter
 *	- r: The http.Request
 *
 * Returns:
 *	- None
 */
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	// Checks the status of the Countries Now API.
	countriesNowStatus := utils.CheckAPIStatus(countriesNowURL)

	// Checks the status of the REST countries API.
	RESTCountriesStatus := utils.CheckAPIStatus(RESTCountriesURL)

	// Calculates the uptime of the service
	uptime := time.Since(startTime).Seconds()

	// Prepare the response
	response := utils.StatusResponse{
		CountriesNowAPI:  countriesNowStatus,
		RestCountriesAPI: RESTCountriesStatus,
		Version:          utils.VERSION,
		Uptime:           uptime,
	}

	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Initiate and encode the response
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error during encoding", http.StatusInternalServerError)
		return
	}

	// Return status code
	http.Error(w, "OK", http.StatusOK)
}
