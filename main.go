package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/icodealot/noaa"
)

// 1. Create one endpoint
//   - Accept latitude and longitude
//   - Returns the short forecast for that area for Today (“Partly Cloudy” etc)
//   - Returns a characterization of whether the temperature is “hot”, “cold”, or “moderate” (use your discretion on mapping temperatures to each type)
//   - Use the National Weather Service API Web Service as a data source.
//
// Example usage: http://localhost:8080/?lat=41.837&lon=-87.685
func main() {
	// Uses standard library http server for simplicity.
	// Production system might use mux, Echo, or Gin frameworks
	// for improved performance, middleware features, or ease or maintenance.

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the route and handler
	mux.Handle("/", &weatherHandler{})

	// Run a simple web server.
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type weatherHandler struct{}

// ServeHTTP provides a single GET endpoint for weather.
//
// A production system would have more multiple routing paths and middleware.
func (h *weatherHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Validate GET verb and single slash routing path.
	if r.Method != http.MethodGet || r.URL.Path != "/" {
		writeJSONError(w, "Not found", http.StatusNotFound)
		return
	}

	qryPrms := r.URL.Query()

	// Validate and parse latitude
	const latLbl = "lat"
	latStr := qryPrms.Get(latLbl)
	if latStr == "" {
		writeJSONError(w, fmt.Sprintf("Query parameter `%v` is missing", latLbl), http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeJSONError(w, fmt.Sprintf("Unable to parse query parameter `%v` to float64", latLbl), http.StatusBadRequest)
		return
	}

	// Validate and parse longitude
	const lonLbl = "lon"
	lonStr := qryPrms.Get(lonLbl)
	if lonStr == "" {
		writeJSONError(w, fmt.Sprintf("Query parameter `%v` is missing", lonLbl), http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		writeJSONError(w, fmt.Sprintf("Unable to parse query parameter `%v` to float64", lonLbl), http.StatusBadRequest)
		return
	}

	// Log the external API call
	log.Printf("External weather request (lat: %.4f, lon: %.4f)", lat, lon)

	// Production system may use context.WithTimeout()
	// for external API call.

	// Call weather API endpoint.
	// Use existing NOAA Go SDK to get weather.
	forecast, err := noaa.Forecast(qryPrms.Get(latLbl), qryPrms.Get(lonLbl))
	if err != nil {
		log.Printf("External weather service error: %v\n", err)
		writeJSONError(w, "External weather service error", http.StatusInternalServerError)
		return
	}
	if len(forecast.Periods) == 0 {
		log.Printf("No weather data returned\n")
		writeJSONError(w, "No weather data returned from external service", http.StatusInternalServerError)
		return
	}

	// Use first period as most relevant forecast.
	p := forecast.Periods[0]

	// Transform data for local HTTP server result.
	res := WeatherResponse{tempLabel(p.Temperature), p.Summary}

	// Serialize to JSON
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		writeJSONError(w, "Error serializing JSON", http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// tempLabel transforms a numeric temperature to a string label.
func tempLabel(temp float64) string {
	// Temperature thresholds in Fahrenheit
	const coldTemp = 40 // Below or equal to this is cold
	const hotTemp = 85  // Above or equal to this is hot

	// String labels
	const coldLabel = "Cold"
	const modLabel = "Moderate"
	const hotLabel = "Hot"

	if temp <= coldTemp {
		return coldLabel
	} else if temp < hotTemp {
		return modLabel
	} else {
		return hotLabel
	}
}

// writeJSONError writes an error in JSON form.
func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResp := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(errResp)
}

// ErrorResponse is an error message for JSON formatting.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WeatherResponse is a weather message for JSON formatting.
type WeatherResponse struct {
	Temperature string `json:"temperature"`
	Forecast    string `json:"forecast"`
}
