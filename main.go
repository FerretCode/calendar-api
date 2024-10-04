package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ferretcode/calendar-api/calendar"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}

	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(5, time.Minute))

	r.Post("/calendar", func(w http.ResponseWriter, r *http.Request) {
		var events []calendar.Event

		err := json.NewDecoder(r.Body).Decode(&events)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "there was an error decoding your events. please make sure they are in a correctly formatted json list.", http.StatusBadRequest)
			return
		}

		img, err := calendar.GenerateCalendar(events)
		if err != nil {
			handleError(err, w)
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(img)
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func handleError(err error, w http.ResponseWriter) {
	fmt.Println(err)
	http.Error(w, "there was an error generating the calendar", http.StatusInternalServerError)
}
