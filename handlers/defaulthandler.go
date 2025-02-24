package handlers

import (
	"assignment1/utils"
	"fmt"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {

	// Ensures client interperates as HTML
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
