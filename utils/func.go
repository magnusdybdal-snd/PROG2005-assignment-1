package utils

import "net/http"

/*
 *	Does a get request to the provided URL, and returns the status code of the response.
 *	Used in checking the status of REST Countries and Countries Now API in statushandler.go
 *
 *	@param - url - The url to invoke and get a response from
 *	@return - The status code of the response
 */
func CheckAPIStatus(url string) int {
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
