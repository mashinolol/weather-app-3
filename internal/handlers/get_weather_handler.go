package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"weather-app-folders/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWeatherHandler(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var weather models.WeatherData
	err := collection.FindOne(ctx, bson.M{"city": city}).Decode(&weather)
	if err != nil {
		http.Error(w, "Weather data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
