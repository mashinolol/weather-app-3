package main

import (
	"fmt"
	"log"
	"net/http"
	"weather-app-3/config"
	"weather-app-3/internal/handlers"
	"weather-app-3/pkg/database"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Подключаемся к MongoDB
	client, collection := database.InitMongoDB(cfg.MongoURI, "weatherdb", "weather")
	defer client.Disconnect(nil)

	// Инициализация HTTP обработчиков
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetWeatherHandler(w, r, collection)
		case http.MethodPut:
			handlers.PutWeatherHandler(w, r, cfg.BaseURL, cfg.APIKey, collection)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
