package handlers

import (
	"assignment1/utils" // const.go, struct.go
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
*	Struct used to populate with the data from the API REST countries, later combined with CityInfo to make our response
 */
type CountryInfo struct {
    Name struct { Common string  `json:"common"`
    } 							 `json:"name"`
    Continents []string          `json:"continents"`
    Population int               `json:"population"`
    Languages  map[string]string `json:"languages"`
    Borders    []string          `json:"borders"`
    Flags struct { PNG string 	 `json:"png"`           // Use the PNG flag URL
    } 							 `json:"flags"`
    Capital []string 			 `json:"capital"`
}

/*
*	Struct used to populate with the data (cities) from the API Countries Now, later combined with CountryInfo to make our response
*/
type CityInfo struct {
	Cities []string `json:"data"`
}

/*
*	Handler that takes care of requests to "/countryinfo/v1/info/" + countryCode. The calls to external API's used in the response
*	is taken care of two seperate functions. Hanlder validates user inputs (country code and limit query) and puts together the response
*	while the seperate functions takes care of invoking the other APIs, and returns those responses to the handler.
*
*	@param - w - Response writer - response sent back to the user
*	@param - r - Request - incoming request from the user
*
*	@see - fetchCountryInfo(string)
*	@see - fetchCities(string)
 */
func InfoHandler(w http.ResponseWriter, r *http.Request) {

	// Extracts the countrycode. using URL.Path will not add optional querries after "?"
	countryCode := strings.TrimPrefix(r.URL.Path, utils.INFO_PATH)

	// Checks if the country code is of len 2, if not a valid code, throw an error
	if len(countryCode) != 2 {
		http.Error(w, "Country code must be a 2-letter ISO code. The code you used is " + strconv.Itoa(len(countryCode)) + " long (\"" + countryCode + "\")", http.StatusBadRequest)
		return
	}

	// Fetch country information from the REST Countries API and checks for error
	countryInfo, err := fetchCountryInfo(countryCode)
	if err != nil {
		log.Printf("Error fetching country information: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch country information.", http.StatusInternalServerError)
		return
	}

	// Fetch city information from the Countries Now API and checks for error
	cityInfo, err := fetchCities(countryInfo.Name.Common)
	if err != nil {
		log.Printf("Error fetching city information: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch city information.", http.StatusInternalServerError)
		return
	}

	limitStr := ""		// used to extract "limit" if present in optional querry (?limit=xx)
	var limit int		// used to extract the number after "=" in optional querry (?limit=xx)

	// Decompose query parameters if present
	if len(r.URL.RawQuery) != 0 {	// checks if there is an optional query

		// Makes sure there are no more than one parameter given.
		if len(r.URL.Query()) != 1 {
			http.Error(w, "Invalid query, only 'limit' parameter is allowed. use (?limit={integer})", http.StatusBadRequest)
			return
		}

		// checks if query is "limit" and extracts the value as string
		limitStr = r.URL.Query().Get("limit")
		// r.URL.Query().Get("limit") returns "" if it cant find "limit" or no int is provided after "="
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)	// Converts the value from string to int
			if err != nil || limit < 0 {		// Checks that limit value is positive
				http.Error(w, "Invalid limit parameter, must be positive integer. (?limit={integer})", http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Invalid query. Use (?limit={integer})", http.StatusBadRequest)
			return
		}
	}

	// Apply the optional limit to cities
	if limitStr != "" && limit < len(cityInfo.Cities){
		cityInfo.Cities = cityInfo.Cities[:limit]
	}

	 // Prepare the response
	 response := utils.InfoResponse{
        Name:       countryInfo.Name.Common,
        Continents: countryInfo.Continents,
        Population: countryInfo.Population,
        Languages:  countryInfo.Languages,
        Borders:    countryInfo.Borders,
        Flag:       countryInfo.Flags.PNG,
        Capital:    countryInfo.Capital[0],
        Cities:     cityInfo.Cities,
    }

    // Set the response content type to JSON
    w.Header().Set("Content-Type", "application/json")

    // Encode the response as JSON and send it
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
    }

	// Return status code
	http.Error(w, "OK", http.StatusOK)
}

/*
*	Calls the REST countries API and retrieves the information needed to fill the CountryInfo struct. 
*
*	@param countryCode - The countrycode used to invoke the REST Countries API (country we are interested in)
*	@return CountryInfo - struct containing the data we are interested in
*	@return - error - error returned if there are issues in fetching and encoding
*/
func fetchCountryInfo(countryCode string) (CountryInfo, error) {

	// URL to invoke
    url := "http://129.241.150.113:8080/v3.1/alpha/" + countryCode

	// Creates new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return CountryInfo{}, fmt.Errorf("error in creating request: %v", err)
	}

	// Sets header
	r.Header.Add("content-type", "application/json")

	// Initiate the client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return CountryInfo{}, fmt.Errorf("error in response: %v", err)
	}
	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		return CountryInfo{}, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	// REST returns a list of countries even tho we ask for only one, we therefore create a list here
	var countryInfo []CountryInfo
	// Decoding json response into a slice of CountryInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&countryInfo); err != nil {
		return CountryInfo{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	//Check if the slice is empty
	if len(countryInfo) == 0 {
		return CountryInfo{}, errors.New("country not found, check if your country code is valid")
	}

	return countryInfo[0], err
}

/*
*	Gets all cities from a specific country countryname provided in the response of the REST country API
*/
func fetchCities(countryName string) (CityInfo, error) {

	// URL to invoke
    url := "http://129.241.150.113:3500/api/v0.1/countries/cities"

	// Payload to use with POST request in countries now API
    payload := map[string]string{"country": countryName}

	// Encodes the payload into JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return CityInfo{}, err
    }

	// Creates new request
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return CityInfo{}, fmt.Errorf("error in creating request: %v", err)
	}

	// Sets header
	r.Header.Add("content-type", "application/json")

	// Initiate the client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return CityInfo{}, fmt.Errorf("error in response: %v", err)
	}
	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		return CityInfo{}, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	// Decodes the JSON response into CityInfo
	var cityInfo CityInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&cityInfo); err != nil {
		return CityInfo{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return cityInfo, err
}