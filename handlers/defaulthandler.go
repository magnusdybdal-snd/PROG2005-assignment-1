package handlers

import (
	"assignment1/utils"
	"fmt"
	"net/http"
)

/*
 * 	DefaultHandler handles requests to the root endpoint (/).
 * 	It provides a simple HTML response with links to other endpoints in the service.
 *
 * 	@param w - The http.ResponseWriter to write the response.
 * 	@param r - The http.Request representing the incoming request.
 */
func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	// Set the response content type to HTML.
	w.Header().Set("content-type", "text/html")

	// Generic info on how to use API and it's endpoints
	output := "No functionality in root level of this service. Please use the following " +
		"paths: <a href=\"" + utils.INFO_PATH + "\">" + utils.INFO_PATH + "</a>, <a href=\"" +
		utils.POPULATION_PATH + "\">" + utils.POPULATION_PATH + "</a> or <a href=\"" + utils.STATUS_PATH +
		"\">" + utils.STATUS_PATH + "</a>."

	// writes output
	_, err := fmt.Fprintf(w, "%v", output)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}

}
