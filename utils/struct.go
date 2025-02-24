package utils

/*
 *	Struct populated with information when /status endpoint is called.
 *	Response is then encoded into json format in statushandler.go
 */
type StatusResponse struct {
	CountriesNowAPI  int     `json:"countriesnowapi"`
	RestCountriesAPI int     `json:"restcountriesapi"`
	Version          string  `json:"version"`
	Uptime           float64 `json:"uptime"`
}

/*
 *	Struct populated with information when /info endpoint is called.
 *	Response is then encoded into json format in infohandler.go
 */
type InfoResponse struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}

/*
 *	Struct populated with information when /population endpoint is called.
 *	Response is then encoded into json format in populationhandler.go
 */
type PopulationResponse struct {
	Mean		int				`json:"mean"`
	Values []struct {
		Year  int				`json:"year"`
		Value int				`json:"value"`
	}							`json:"values"`
}
