package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rcbadiale/go-cloud-run/internals/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file, will use environment variables")
	}

	log.Println("starting server on port 8080")
	runServer()
}

func runServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	weatherApiKey := os.Getenv("WEATHER_API_KEY")
	weatherHandler := handlers.NewWeatherHandler(weatherApiKey)
	r.With(addContext).Get("/weather/{zipCode}", weatherHandler.GetWeather)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Println("error starting server: ", err)
	}
}

type contextKey string

func addContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey("request_id"), uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
