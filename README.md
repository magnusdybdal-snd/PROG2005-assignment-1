# Country Information Service

This is a Go-based web service that provides information about countries, including population data, cities, languages, and more. It integrates with external APIs [REST Countries](http://129.241.150.113:8080/v3.1/)and [CountriesNow](http://129.241.150.113:3500/api/v0.1/) (self hosted) to fetch and serve data.

The service is deployed on **Render** and can be accessed via the provided Render URL. The source code is hosted on **GitLab** and **GitHub** and is available for internal viewing.

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

- **Country Information**: Retrieve detailed information about a country, including its name, continents, population, languages, borders, flag, capital and cities.
- **Population Data**: Fetch population statistics for a country, optionally filtered by a year range.
- **Service Status**: Check the status of the service and external APIs.

---

## Endpoints

### 1. **Root Endpoint**
- **Path**: `/`
- **Method**: `GET`
- **Description**: Provides a short message and links to other endpoints.
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

### Installation
1. Clone the repository:
   ```bash
   git clone https://gitlab.com/your-username/your-repo-name.git
   cd your-repo-name

2. Build the project:
    ```bash
    go build

3. Run the service:
    ```bash
    ./your-repo-name

---

## Usage

### Example Requests

1. **Get Country Information**:
   ```bash
   curl http://localhost:8080/countryinfo/v1/info/no

### Response

    ```json
    {
    "name": "Norway",
    "continents": ["Europe"],
    "population": 5379475,
    "languages": {
        "nor": "Norwegian"
    },
    "borders": ["FIN", "SWE", "RUS"],
    "flag": "https://flagcdn.com/no.svg",
    "capital": "Oslo",
    "cities": ["Oslo", "Bergen", "Trondheim"]
    }