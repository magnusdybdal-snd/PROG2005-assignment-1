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
	output := "<h2>Welcome to the Prog2005 - assignment1 service root page.</h2>" + 
	    " <p>Please use the following paths to access the service:</p> "     +
		" <a href=\"" + utils.INFO_PATH       + "\">" + utils.INFO_PATH       + "</a><code>{2 letter ISO code}</code>,"   +
		" <a href=\"" + utils.POPULATION_PATH + "\">" + utils.POPULATION_PATH + "</a><code>{2 letter ISO code}</code> or" +
		" <a href=\"" + utils.STATUS_PATH     + "\">" + utils.STATUS_PATH     + "</a>." +
		" <p>Replace <code>{2 letter ISO code}</code> with the appropriate country code (e.g.  <code>\"no\"</code>  for Norway)" +
		" <h3>Optional queries</h3>" +
		" <p>Info/<code>{2 letter ISO code}</code> has an optional query: <code>?limit=X</code> where <code>X</code> is an integer to limit the amount of cities displayed</p>"  +
		" <p>Population/<code>{2 letter ISO code}</code> has an optional query <code>?limit=X-Y</code> where <code>X</code> and <code>Y</code> are integers, to filter population data</p>"  +
		" <h3>Example endpoints:</h3>" +
		" <p><a href=\"" + utils.INFO_PATH + "no?limit=10" + "\">" + utils.INFO_PATH + "no?limit=10" + "</a></p>" +
		" <p><a href=\"" + utils.POPULATION_PATH + "no?limit=2000-2005" + "\">" + utils.POPULATION_PATH + "no?limit=2000-2005" + "</a></p>" +
		" <p><a href=\"" + utils.STATUS_PATH  + "\">" + utils.STATUS_PATH + "</a></p>"

	// writes output
	_, err := fmt.Fprintf(w, "%v", output)
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}

}
