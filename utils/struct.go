package utils

/*
*	The first 3 structs with names ending in Response are used to construct and represent the response
* 	for the endpoints of this service, before they are encoded into json format.
* 	The last 4 structs with names ending in Info are used to fetch information from REST countries and
* 	CountriesNow API where their json responses are decoded into the structs.
 */

/*
 * 	StatusResponse represents the response structure for the /status endpoint.
 * 	It includes the status of external APIs, the service version, and uptime.
 */
 type StatusResponse struct {
	CountriesNowAPI  int     `json:"countriesnowapi"`  	// Status code for the CountriesNow API
	RestCountriesAPI int     `json:"restcountriesapi"` 	// Status code for the REST Countries API
	Version          string  `json:"version"`          	// Service version
	Uptime           float64 `json:"uptime"`           	// Service uptime in seconds
}

/*
 * 	InfoResponse represents the response structure for the /info endpoint.
 * 	It includes detailed information about a country, such as its name, population, and cities.
 */
 type InfoResponse struct {
	Name       string            `json:"name"`       	// Common name of the country
	Continents []string          `json:"continents"` 	// List of continents the country belongs to
	Population int               `json:"population"` 	// Total population of the country
	Languages  map[string]string `json:"languages"`  	// Map of language codes to language names
	Borders    []string          `json:"borders"`    	// List of bordering country codes
	Flag       string            `json:"flag"`       	// URL to the country's flag (PNG format)
	Capital    string            `json:"capital"`    	// Capital city of the country
	Cities     []string          `json:"cities"`     	// List of cities in the country
}

/*
 * 	PopulationResponse represents the response structure for the /population endpoint.
 * 	It includes the mean population value and a list of population data over time.
 */
 type PopulationResponse struct {
	Mean   int 		  `json:"mean"` 					// Mean population value over the specified years
	Values []struct {
		Year  int 	  `json:"year"`  					// Year of the population data
		Value int 	  `json:"value"` 					// Population value for the corresponding year
	} 				  `json:"values"` 					// List of population data points
}

/*
 * 	CountryInfo represents parts of the structure of data returned by the REST Countries API.
 * 	It is used internally to fetch country information for the /info endpoint.
 */
 type CountryInfo struct {
	Name struct { Common string    `json:"common"` 		// Common name of the country
	} 							   `json:"name"`
	Continents []string            `json:"continents"` 	// List of continents the country belongs to
	Population int                 `json:"population"` 	// Total population of the country
	Languages  map[string]string   `json:"languages"`  	// Map of language codes to language names
	Borders    []string            `json:"borders"`    	// List of bordering country codes
	Flags      struct { PNG string `json:"png"` 		// URL to the country's flag (PNG format)
	} 							   `json:"flags"`
	Capital []string 			   `json:"capital"`		// List of capital citie(s)
}

/*
 * CountryName represents a subset of the data returned by the REST Countries API.
 * It is used internally to fetch the common name of a country based on its country code,
 * which is required for querying population data in the /population endpoint.
 */
type CountryName struct {
	Name struct { Common string `json:"common"`			// Common name of the country
   } 							`json:"name"`
}

/*
 * 	CityInfo represents parts of the structure of data returned by the CountriesNow API.
 * 	It is used internally to fetch city information for the /info endpoint.
 */
type CityInfo struct {
	Cities []string `json:"data"` 						// List of cities in the country
}


/*
 * 	PopulationInfo represents the structure of data returned by the CountriesNow API.
 * 	It is used internally to fetch population data for the /population endpoint.
 */
type PopulationInfo struct {
	Data struct {
		PopulationCount []struct {
			Year  int `json:"year"`  					// Year of the population data
			Value int `json:"value"` 					// Population value for the corresponding year
		} 			  `json:"populationCounts"` 		// List of population data points
	} 			      `json:"data"`
}