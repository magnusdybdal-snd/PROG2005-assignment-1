# Country Information Service

This is a Go-based web service that provides information about countries, including population data, cities, languages, and more. It integrates with external APIs like [REST Countries](https://restcountries.com/) and [CountriesNow](https://countriesnow.space/) to fetch and serve data.

The service is deployed on **Render** and can be accessed via the provided Render URL. The source code is hosted on **GitLab** and is available for internal viewing by classmates.

---

## Table of Contents
1. [Features](#features)
2. [Endpoints](#endpoints)
3. [Setup](#setup)
4. [Usage](#usage)
5. [Deployment](#deployment)
6. [API Documentation](#api-documentation)
7. [Contributing](#contributing)
8. [License](#license)

---

## Features

- **Country Information**: Retrieve detailed information about a country, including its name, population, languages, borders, and more.
- **Population Data**: Fetch population statistics for a country, optionally filtered by a year range.
- **Service Status**: Check the status of the service and external APIs.
- **Simple Interface**: Easy-to-use endpoints with clear responses in JSON format.

---

## Endpoints

### 1. **Root Endpoint**
- **Path**: `/`
- **Method**: `GET`
- **Description**: Provides a welcome message and links to other endpoints.
- **Response**: HTML page with links.

### 2. **Country Information**
- **Path**: `/countryinfo/v1/info/{countryCode}`
- **Method**: `GET`
- **Description**: Retrieves detailed information about a country.
- **Query Parameters**:
  - `limit` (optional): Limits the number of cities returned.
- **Response**: JSON object with country details.

### 3. **Population Data**
- **Path**: `/countryinfo/v1/population/{countryCode}`
- **Method**: `GET`
- **Description**: Retrieves population data for a country.
- **Query Parameters**:
  - `limit` (optional): Filters population data by a year range (e.g., `?limit=2000-2020`).
- **Response**: JSON object with population data.

### 4. **Service Status**
- **Path**: `/countryinfo/v1/status`
- **Method**: `GET`
- **Description**: Checks the status of the service and external APIs.
- **Response**: JSON object with service status.

---

## Setup

### Prerequisites
- Go 1.22.2 or higher
- Git (optional, for cloning the repository)

### Installation (locally instead of using the Render service)
1. Clone the repository:
   ```bash
   git clone https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2025-workspace/magnusdybdal/assignment-1.git
    ```
   ```bash
   cd assignment-1

2. Build the project:
    ```bash
   go build

3. Run the service:
    ```bash
   ./assignment1

---

## Usage

### Example Requests

1. **Get Country Information**:
    ```bash
   curl http://localhost:8080/countryinfo/v1/info/no?limit=5

Response:
```json
{
    "name":"Norway",
    "continents": [
        "Europe"
    ],
    "population":5379475,
    "languages": {
        "nno":"Norwegian Nynorsk",
        "nob":"Norwegian Bokmål",
        "smi":"Sami"
    },"borders": [
        "FIN",
        "SWE",
        "RUS"
    ],
    "flag":"https://flagcdn.com/w320/no.png",
    "capital":"Oslo",
    "cities": [
        "Abelvaer",
        "Adalsbruk",
        "Adland",
        "Agotnes",
        "Agskardet"
    ]
}
```

2. **Get Population Data**:
    ```bash
   curl http://localhost:8080/countryinfo/v1/population/no?limit=2000-2002

Response:
```json
   {
     "mean": 4514292,
     "values": [
       {
           "year": 2000, 
           "value": 4490967
       },
       {
           "year": 2001, 
           "value": 4513751
       },
       {
           "year": 2002, 
           "value": 4538159
       },
     ]
   }
```

3. **Check Service Status**:
    ```bash
   curl http://localhost:8080/countryinfo/v1/status

Response:
    
```json
{
    "countriesnowapi": 200,
    "restcountriesapi": 200,
    "version": "v1",
    "uptime": 3852.541554949
}
```
---

## Deployment

The service is deployed on **Render**. You can access it at the following URL:

🔗 **[Render Service URL](https://assignment-1-uszz.onrender.com/)**


---