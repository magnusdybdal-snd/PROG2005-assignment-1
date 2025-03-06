package utils

// DEFAULT_PATH is the root endpoint of the service.
const DEFAULT_PATH = "/"

// INFO_PATH is the endpoint for retrieving country information.
const INFO_PATH = "/countryinfo/" + VERSION + "/info/"

// POPULATION_PATH is the endpoint for retrieving population data.
const POPULATION_PATH = "/countryinfo/" + VERSION + "/population/"

// STATUS_PATH is the endpoint for retrieving service status.
const STATUS_PATH = "/countryinfo/" + VERSION + "/status"

// VERSION represents the current version of the service.
const VERSION = "v1"

// DEFAULT_PORT is the default port the service listens on if no port is specified.
const DEFAULT_PORT = "8080"

// RESTCountriesURL is the base URL for the REST Countries API.
const RESTCountriesURL = "http://129.241.150.113:8080/v3.1/alpha/"

// CountriesNowURL is the base URL for the CountriesNow API.
const CountriesNowURL = "http://129.241.150.113:3500/api/v0.1/countries/"