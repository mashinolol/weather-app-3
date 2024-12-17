package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"weather-app-3/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PutWeatherHandler(w http.ResponseWriter, r *http.Request, baseURL, apiKey string, collection *mongo.Collection) {
	var requestBody struct {
		City string `json:"city"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil || requestBody.City == "" {
		http.Error(w, "Invalid or missing city parameter", http.StatusBadRequest)
		return
	}

	// Fetch weather data from OpenWeatherMap API
	searchURL := fmt.Sprintf("%v?appid=%s&q=%s", baseURL, apiKey, requestBody.City)
	resp, err := http.Get(searchURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	weatherBytes, _ := io.ReadAll(resp.Body)
	var weatherAPIResponse struct {
		Weather []struct{ Description string } `json:"weather"`
		Main    struct{ Temp float64 }         `json:"main"`
		Name    string                         `json:"name"`
	}
	if err := json.Unmarshal(weatherBytes, &weatherAPIResponse); err != nil {
		http.Error(w, "Failed to parse weather data", http.StatusInternalServerError)
		return
	}

	weatherData := models.WeatherData{
		City:        weatherAPIResponse.Name,
		Description: weatherAPIResponse.Weather[0].Description,
		Temp:        weatherAPIResponse.Main.Temp - 273.15,
		LastUpdated: time.Now(),
	}

	// Upsert to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"city": weatherData.City}
	update := bson.M{"$set": weatherData}
	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		http.Error(w, "Failed to update weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weatherData)
}
