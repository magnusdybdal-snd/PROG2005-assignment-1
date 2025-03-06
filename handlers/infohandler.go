package handlers

import (
	"assignment1/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
 * 	InfoHandler handles requests to the /info endpoint.
 * 	It retrieves detailed information about a country, including its name, population, languages, and cities using 2 API's_
 *	REST Countries and CountriesNow
 *
 * 	@param w - The http.ResponseWriter to write the response.
 * 	@param r - The http.Request representing the incoming request.
 */
func InfoHandler(w http.ResponseWriter, r *http.Request) {

	// Method only allows GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extracts the countrycode. using URL.Path ignores queries when extracting
	countryCode := strings.TrimPrefix(r.URL.Path, utils.INFO_PATH)

	// Validate the country code length
	if len(countryCode) != 2 {
		http.Error(w, "Country code must be a 2-letter ISO code. for example try: " + utils.INFO_PATH + "no?limit=10", http.StatusBadRequest)
		return
	}

	// Fetch country information from the REST Countries API
	countryInfo, err := fetchCountryInfo(countryCode)
	if err != nil {
		log.Printf("Error fetching country information: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch country information.", http.StatusInternalServerError)
		return
	}

	// Fetch city information from the Countries Now API
	cityInfo, err := fetchCities(countryInfo.Name.Common)
	if err != nil {
		log.Printf("Error fetching city information: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch city information.", http.StatusInternalServerError)
		return
	}
		
		// checks if query is "limit" and extracts the value as string
		limitStr := r.URL.Query().Get("limit")
		var limit int		// used to convert limit string into an int
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)	// Converts the value from string to int
			if err != nil || limit < 0 {		// Checks that limit value is positive
				http.Error(w, "Invalid limit parameter, must be positive integer.", http.StatusBadRequest)
				return
			}
			// Apply the optional limit to the list of cities
			if limit < len(cityInfo.Cities){
				cityInfo.Cities = cityInfo.Cities[:limit]
			}
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

	// Encode the response as JSON and send it. This also sets the status code to 200 if successful
    if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v\n", err)
        http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
    }
}

/*
 * 	fetchCountryInfo retrieves country information from the REST Countries API using a country code.
 *
 * 	@param countryCode - The 2-letter ISO country code.
 * 	@return utils.CountryInfo - The country information.
 * 	@return error - An error if the request or decoding fails.
 */
func fetchCountryInfo(countryCode string) (utils.CountryInfo, error) {

	// URL to invoke
    url := utils.RESTCountriesURL + countryCode

	// Uses http.Get because we don't need custom headers, client or timeouts. The header content type is automatically set to JSON with http.Get
	r, err := http.Get(url) 
	if err != nil {
		return utils.CountryInfo{}, fmt.Errorf("error fetching country info: %v", err)
	}
	defer r.Body.Close()

	// Check the HTTP status code
	if r.StatusCode != http.StatusOK {
		return utils.CountryInfo{}, fmt.Errorf("API returned non-200 status code: %d", r.StatusCode)
	}

	// REST returns a list of countries even tho we ask for only one, we therefore create a list here
	var countryInfo []utils.CountryInfo
	// Decoding json response into a slice of CountryInfo
	if err := json.NewDecoder(r.Body).Decode(&countryInfo); err != nil {
		return utils.CountryInfo{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	//Check if the slice is empty
	if len(countryInfo) == 0 {
		return utils.CountryInfo{}, fmt.Errorf("countryInfo slice is empty: %v", err)
	}

	return countryInfo[0], nil
}

/*
 * 	fetchCities retrieves city information from the CountriesNow API using a country name.
 *
 * 	@param countryName - The name of the country.
 * 	@return utils.CityInfo - The list of cities in the country.
 * 	@return error - An error if the request or decoding fails.
 */
func fetchCities(countryName string) (utils.CityInfo, error) {

	// URL to invoke
    url := utils.CountriesNowURL + "cities"

	// Creates payload to use with POST request in countries now API and encodes it into JSON
    payload := map[string]string{"country": countryName}
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return utils.CityInfo{}, fmt.Errorf("failed to encode payload: %v", err)
    }

	// Creates new request
	r, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return utils.CityInfo{}, fmt.Errorf("failed to fetch city information: %v", err)
	}
	defer r.Body.Close()

	// Check the HTTP status code
	if r.StatusCode != http.StatusOK {
		return utils.CityInfo{}, fmt.Errorf("API returned non-200 status code: %d", r.StatusCode)
	}

	// Decodes the JSON response into CityInfo
	var cityInfo utils.CityInfo
	if err := json.NewDecoder(r.Body).Decode(&cityInfo); err != nil {
		return utils.CityInfo{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return cityInfo, nil
}