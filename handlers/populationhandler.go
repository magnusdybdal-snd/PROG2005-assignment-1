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
 * 	PopulationHandler handles requests to the /population endpoint.
 * 	It retrieves population data for a specific country and optionally filters it by year range.
 *
 * 	@param w - The http.ResponseWriter to write the response.
 * 	@param r - The http.Request representing the incoming request.
 */
func PopulationHandler(w http.ResponseWriter, r *http.Request) {

	// Endpoint only allows GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the country code from the URL path.
	countryCode := strings.TrimPrefix(r.URL.Path, utils.POPULATION_PATH)

	// Validate the country code length.
	if len(countryCode) != 2 {
		http.Error(w, "Country code must be a 2-letter ISO code. for example try: " + utils.POPULATION_PATH + "no?limit=2000-2010", http.StatusBadRequest)
		return
	}

	// Fetch country name  from the REST Countries API using the country code.
	countryName, err := fetchCountryName(countryCode)
	if err != nil {
		log.Printf("error fetching country name: %v\n", err)	// Log specific error
		http.Error(w, "failed to fetch country name", http.StatusInternalServerError)
		return
	}

	// Fetch population information from the Countries now API
	populationInfo, err := fetchPopulation(countryName.Name.Common)
	if err != nil {
		log.Printf("error fetching population data: %v\n", err)	// Log specific error
		http.Error(w, "failed to fetch population information", http.StatusInternalServerError)
		return
	}

	// start and end year of optional limit query
	var startYear, endYear int		
	// Decompose query parameters if present
	if len(r.URL.RawQuery) != 0 {		// checks if optional query is present
		
		// checks if query is "limit" and has a valid value
		limitStr := r.URL.Query().Get("limit")
		if limitStr != "" {

			years := strings.Split(limitStr, "-")
			if len(years) != 2 {		// Checks if the user entered 2 year split by -
				http.Error(w, "invalid format for limit. use (?limit=startYear-endYear)", http.StatusBadRequest)
				return
			}
			// Ensures that startYear is a valid integer
			var err error
			startYear, err = strconv.Atoi(years[0])
			if err != nil {
				http.Error(w, "invalid start year, must be an integer.", http.StatusBadRequest)
				return
			}

			// Ensures that endYear is a valid integer
			endYear, err = strconv.Atoi(years[1])
			if err != nil {
				http.Error(w, "invalid end year, must be an integer.", http.StatusBadRequest)
				return
			}

			// Validate the year range
			if startYear > endYear {
				http.Error(w, "endyear must be greater than or equal to startyear", http.StatusBadRequest)
				return
			}
			
		} else {	
			http.Error(w, "invalid query. Use (?limit=integer-integer)", http.StatusBadRequest) // If query is not "limit"
			return
		}
	}

	// Filter population data by year range if applicable.
	var filteredPopulation []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	}
	var meanValue int

	if startYear != 0 || endYear != 0 {
		// Loops trough the results from CountriesNow and adds the years we are interested in to the struct
		for _, entry := range populationInfo.Data.PopulationCount {
			if entry.Year >= startYear && entry.Year <= endYear {
				filteredPopulation = append(filteredPopulation, entry)
			}
		}
		
		// Checks if filter is within scope of the data provided by CountriesNow API
		if len(filteredPopulation) == 0 {
			http.Error(w, "no population data within the span of years provided. first recorded year is " + 
			strconv.Itoa(populationInfo.Data.PopulationCount[0].Year) , http.StatusBadRequest)
			return
		}
		
		// Iterates trough the filtered population data and calculates the mean value
		sum := 0
		for _, year := range filteredPopulation {
			sum += year.Value
		}
		meanValue = sum / len(filteredPopulation)
	} else {
		// Use all population data if no filtering is applied
		filteredPopulation = populationInfo.Data.PopulationCount
	}

	// Preapering the response
	response := utils.PopulationResponse{
		Mean:   meanValue,
		Values: filteredPopulation,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and send it. This also sets the status code to 200 if successful
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("error encoding response: %v\n", err)	// Log specific error
		http.Error(w, "failed to encode response.", http.StatusInternalServerError)
		return
	}
}

/*
 * 	fetchCountryName retrieves the country name from the REST Countries API using a country code.
 *
 * 	@param countryCode - The 2-letter ISO country code.
 * 	@return utils.CountryName - The full country name belonging to the country code
 * 	@return error - An error if the request or decoding fails.
 */
func fetchCountryName (countryCode string) (utils.CountryName, error){
	
	// URL to invoke
    url := utils.RESTCountriesURL + countryCode

	// Uses http.Get because we don't need custom headers, client or timeouts. The header content type is automatically set to JSON with http.Get
	r, err := http.Get(url)
	if err != nil {
		return utils.CountryName{}, fmt.Errorf("failed to fetch country name: %v", err)
	}
	defer r.Body.Close()

	// Check the HTTP status code of the response
	if r.StatusCode != http.StatusOK {
		return utils.CountryName{}, fmt.Errorf("api returned non-200 status code: %d", r.StatusCode)
	}

	// REST returns a list of countries even tho we ask for only one, we therefore create a list here
	var countryName []utils.CountryName
	// Decoding json response into a slice of CountryInfo
	if err := json.NewDecoder(r.Body).Decode(&countryName); err != nil {
		return utils.CountryName{}, fmt.Errorf("failed to decode JSON response: %v", err)
	}

	return countryName[0], nil
}

/*
 * 	fetchPopulation retrieves population data from the CountriesNow API using a country name.
 *
 * 	@param countryName - The name of the country.
 * 	@return utils.PopulationInfo - The population data for the country.
 * 	@return error - An error if the request or decoding fails.
 */
func fetchPopulation (countryName string) (utils.PopulationInfo, error){

	// URL to invoke
    url := utils.CountriesNowURL + "population"

	// Creates payload to use with POST request in countries now API and encodes into JSON
    payload := map[string]string{"country": countryName}
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return utils.PopulationInfo{}, fmt.Errorf("failed to encode payload: %v", err)
    }

	r, err := http.Post(url, "application/json", strings.NewReader(string(jsonPayload)))
	if err != nil {
		return utils.PopulationInfo{}, fmt.Errorf("failed to fetch population data: %v", err)
	}
	defer r.Body.Close()

	// Check the HTTP status code
	if r.StatusCode != http.StatusOK {
		return utils.PopulationInfo{}, fmt.Errorf("api returned non-200 status code: %d", r.StatusCode)
	}

	// Decodes the JSON response into populationInfo
	var populationInfo utils.PopulationInfo
	if err := json.NewDecoder(r.Body).Decode(&populationInfo); err != nil {
		return utils.PopulationInfo{}, fmt.Errorf("error decoding json: %v", err)
	}

	return populationInfo, nil
}