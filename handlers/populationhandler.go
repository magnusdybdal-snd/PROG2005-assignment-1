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
*	TODO KOMMENTARER
 */
func PopulationHandler(w http.ResponseWriter, r *http.Request) {

	// Extracts the countrycode. using URL.Path will not add the optional ?limit=xx
	countryCode := strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/population/")

	// Checks if the country code is of len 2, if not a valid code, the API will throw error
	if len(countryCode) != 2 {
		http.Error(w, "Country code must be a 2-letter ISO code. The code you used is " + strconv.Itoa(len(countryCode)) + " long (\"" + countryCode + "\")", http.StatusBadRequest)
		return
	}

	// Fetch country name from the REST Countries API and checks for error
	countryName, err := fetchCountryName(countryCode)
	if err != nil {
		log.Printf("Error fetching country name: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch country name", http.StatusInternalServerError)
		return
	}

	// Fetch population information from the Countries now API and checks for error
	populationInfo, err := fetchPopulation(countryName.Name.Common)
	if err != nil {
		log.Printf("Error fetching population data: %v\n", err)	// Log specific error
		http.Error(w, "Failed to fetch population information", http.StatusInternalServerError)
		return
	}

	// Optional query {?limit={:startyear-endyear}} and mean value
	var limitStr = ""		// optional limit query as string
	var startYear int		// start year of selection
	var endYear int			// end year of selection
	var meanValue int		// the calculated mean value

	// Decompose query parameters if present
	if len(r.URL.RawQuery) != 0 {		// checks if optional query is present

		// checks if query is "limit" and has a valid value
		limitStr = r.URL.Query().Get("limit")
		if limitStr != "" {
			var err error
			years := strings.Split(limitStr, "-")
			if len(years) != 2 {		// Checks if the user entered 2 year split by -
				http.Error(w, "Invalid format for limit. use (?limit=integer-integer)", http.StatusBadRequest)
				return
			}
			// Ensures that startYear is a valid integer
			startYear, err = strconv.Atoi(years[0])
			if err != nil {
				http.Error(w, "Invalid start year, Must be an integer. use (?limit=integer-integer)", http.StatusBadRequest)
				return
			}

			// Ensures that endYear is a valid integer
			endYear, err = strconv.Atoi(years[1])
			if err != nil {
				http.Error(w, "Invalid end year, must be an integer. use (?limit=integer-integer)", http.StatusBadRequest)
				return
			}

			// Ensures that startYear is >= endYear
			if startYear > endYear {
				http.Error(w, "Endyear must be greater than or equal to startyear", http.StatusBadRequest)
				return
			}

			// Apply the filtering of years chosen
			// Helping struct to hold the values years and values requested
			filteredPopulation := []struct {
				Year  int `json:"year"`
				Value int `json:"value"`
			}{}

			// Loops trough the results from CountriesNow and adds the years we are interested in
			// To the helping struct
			for _, entry := range populationInfo.Data.PopulationCount {
				if entry.Year >= startYear && entry.Year <= endYear{
					filteredPopulation = append(filteredPopulation, entry)
				}
			}
			
			// Checks if filter is within scope of the data provided by CountriesNow
			if len(filteredPopulation) == 0 {
				http.Error(w, "No population data within the span of years provided. First recorded year is " + 
				strconv.Itoa(populationInfo.Data.PopulationCount[0].Year) , http.StatusBadRequest)
				return
			}
			
			// Iterates trough the filtered population data and calculates the mean (as int)
			sum := 0
			for _, year := range filteredPopulation {
				sum += year.Value
			}
			meanValue = sum / len(filteredPopulation)
			// Updates the data returned from CountriesNow with our filter
			populationInfo.Data.PopulationCount = filteredPopulation


		} else {	// If query is not "limit"
			http.Error(w, "Invalid query. Use (?limit=integer-integer)", http.StatusBadRequest)
			return
		}
	}

	// Preapering the response
	response := utils.PopulationResponse{
		Mean: meanValue,
		Values: populationInfo.Data.PopulationCount,
	}

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the response as JSON and send it
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v\n", err)	// Log specific error
		http.Error(w, "Failed to encode response.", http.StatusInternalServerError)
	}

	// Return status code
	http.Error(w, "OK", http.StatusOK)
}

/*
*	TODO KOMMENTARER (Flytt struct?)
*/
type CountryName struct {
	 Name struct { Common string  `json:"common"`
    } 							  `json:"name"`
}

/*
*	Uses the REST country API to fetch the country name, using a country code
*/
func fetchCountryName (countryCode string) (CountryName, error){
	
	// URL to invoke
    url := "http://129.241.150.113:8080/v3.1/alpha/" + countryCode

	// Creates new request
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return CountryName{}, fmt.Errorf("error in creating request: %v", err)
	}

	// Sets header
	r.Header.Add("content-type", "application/json")

	// Initiate the client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return CountryName{}, fmt.Errorf("error in response: %v", err)
	}
	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		return CountryName{}, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	// REST returns an array of countries even tho we ask for only one, we therefore create an array here
	var countryName []CountryName
	// Decoding json response into a slice of CountryInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&countryName); err != nil {
		return CountryName{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	return countryName[0], err
}

/*
*	TODO KOMMENTARER
*/
type PopulationInfo struct {
	Data struct { 
		PopulationCount []struct{
			Year int  	`json:"year"`
			Value int 	`json:"value"`
		}				`json:"populationCounts"`
	}					`json:"data"`
}

/*
*	TODO KOMMENTARER
*/
func fetchPopulation (countryName string) (PopulationInfo, error){
	// URL to invoke
    url := "http://129.241.150.113:3500/api/v0.1/countries/population"

	// Payload to use with POST request in countries now API
    payload := map[string]string{"country": countryName}

	// Encodes the payload into JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return PopulationInfo{}, err
    }

	// Creates new request
	r, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return PopulationInfo{}, fmt.Errorf("error in creating request: %v", err)
	}

	// Sets header
	r.Header.Add("content-type", "application/json")

	// Initiate the client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// Issue request
	res, err := client.Do(r)
	if err != nil {
		return PopulationInfo{}, fmt.Errorf("error in response: %v", err)
	}
	defer res.Body.Close()

	// Check the HTTP status code
	if res.StatusCode != http.StatusOK {
		return PopulationInfo{}, fmt.Errorf("API returned non-200 status code: %d", res.StatusCode)
	}

	// Decodes the JSON response into populationInfo
	var populationInfo PopulationInfo
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&populationInfo); err != nil {
		return PopulationInfo{}, fmt.Errorf("error decoding JSON: %v", err)
	}

	fmt.Println(populationInfo)
	return populationInfo, err
}