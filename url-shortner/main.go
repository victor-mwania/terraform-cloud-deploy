package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/speps/go-hashids"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type ShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

type URL struct {
	gorm.Model
	ShortUrl string `json:"short_url"`
	FullUrl  string `json:"full_url"`
}

func main() {
	dsn := "host=localhost user=postgres password=password dbname=db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&URL{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")

	r := chi.NewRouter()

	r.Get("/{shortURL}", func(w http.ResponseWriter, r *http.Request) {
		shortURL := chi.URLParam(r, "shortURL")

		fmt.Println(shortURL, "shortURL")
		var url URL
		result := db.Where(&URL{ShortUrl: shortURL}).First(&url)
		// result := db.First(&url)
		fmt.Println(url.FullUrl, "result")
		if result.Error != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, url.FullUrl, http.StatusSeeOther)
	})

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		fullURL := r.FormValue("url")

		// Generate a short ID using hashids.
		hd := hashids.NewData()
		hd.Salt = "url_shortner_salt" // should be read from config
		h, _ := hashids.NewWithData(hd)
		shortID, _ := h.Encode([]int{123})

		// Store the shortID and fullURL in the database.
		url := URL{ShortUrl: shortID, FullUrl: fullURL}
		result := db.Create(&url)
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := ErrorResponse{Message: "Error storing URL"}
			jsonResponse, _ := json.Marshal(errorResponse)
			_, _ = w.Write(jsonResponse)
			return
		}

		shortURL := fmt.Sprintf("/%s", shortID)
		response := ShortURLResponse{ShortURL: shortURL}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(jsonResponse)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := []byte(`{"status": "ok"}`)
		_, _ = w.Write(response)
	})

	port := ":8090"
	log.Printf("Server is listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}
