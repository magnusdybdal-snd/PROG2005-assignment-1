package handlers

import (
	"assignment1/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// startTime tracks when the service started, used to calculate uptime.
var startTime = time.Now()

/*
 * 	StatusHandler handles requests to the /status endpoint.
 * 	It checks the status of external APIs (CountriesNow and REST Countries),
 * 	calculates the service uptime, and returns a JSON response.
 *
 * 	@param w - The http.ResponseWriter to write the response.
 * 	@param r - The http.Request representing the incoming request.
 */
func StatusHandler(w http.ResponseWriter, r *http.Request) {

	// Method only allows GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Checks the status of the Countries Now API.
	countriesNowStatus := CheckAPIStatus(utils.CountriesNowURL + "iso")

	// Checks the status of the REST countries API.
	RESTCountriesStatus := CheckAPIStatus(utils.RESTCountriesURL + "no")

	// Calculates the uptime of the service in seconds
	uptime := time.Since(startTime).Seconds()

	// Prepare the response
	response := utils.StatusResponse{
		CountriesNowAPI:  countriesNowStatus,
		RestCountriesAPI: RESTCountriesStatus,
		Version:          utils.VERSION,
		Uptime:           uptime,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and send it. if successful will also add status 200 to header
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log the error and return a 500 Internal Server Error.
		log.Printf("failed to encode response: %v\n", err)
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
		return
	}
}

/*
 * 	CheckAPIStatus sends a GET request to the specified URL and returns the HTTP status code.
 * 	If the request fails, it returns http.StatusServiceUnavailable (503).
 *
 * 	@param url - The URL to send the GET request to.
 * 	@return - The HTTP status code of the response.
 */
 func CheckAPIStatus(url string) int {
	// Validate the URL
	if url == "" {
		return http.StatusServiceUnavailable
	}

	// Sending a GET request to the url
	resp, err := http.Get(url)
	// If request fails we return 503
	if err != nil {
		return http.StatusServiceUnavailable
	}
	// Ensures the response body is closed
	defer resp.Body.Close()

	// Return the status code
	return resp.StatusCode
}
