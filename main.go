package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome"))
	})

	router.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
		city := r.URL.Query().Get("city")
		latLong, err := getLatLong(city)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to get Latitude and Longitude %s", err)))
			return
		}

		weather, err := getWeather(*latLong)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to get weather %s", err)))
			return
		}

		fmt.Println(weather)

		weatherDisplay, err := extractWeatherData(city, weather)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("failed to extract weather data %s", err)))
			return
		}
		w.Write([]byte(fmt.Sprintf("weather: %s", weather)))
		t, _ := template.ParseFiles("views/weather.html")
		t.Execute(w, weatherDisplay)

	})

	fmt.Printf("Server running on Port:%s\n", port)
	http.ListenAndServe(":"+port, router)
}
