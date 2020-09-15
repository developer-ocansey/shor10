package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"urlshortner/models"
	"urlshortner/utils"

	"github.com/gorilla/mux"
)

// Response ...
type Response struct {
	status     bool
	shortenURL string
	error      string
}

// Value ...
type Value struct {
	id           int
	generatedURL string
	originalURL  string
}

// Query ...
type Query struct {
	Value        Value
	Error        string
	RowsAffected int
}

// URL ...
type URL struct {
	URL string `json:"url"`
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	var bodyURL URL

	if err := json.NewDecoder(r.Body).Decode(&bodyURL); err != nil {
		log.Fatal(err)
	}

	url := bodyURL.URL
	if len(url) < 25 {
		utils.FormatResponse("invalid", "URL must be more than 24 characters", w)
	} else {
		hash := sha1.New()
		hash.Write([]byte(url))
		hashedURL := string(hash.Sum(nil))
		sliceURL := hashedURL[0:3]
		shortURL := fmt.Sprintf("%x", sliceURL)

		data := models.StoreURL{GeneratedURL: shortURL, OriginalURL: url}
		response := db.Where(&models.StoreURL{OriginalURL: url}).Find(&data)

		if response.RowsAffected > 0 {
			utils.FormatResponse("exist", "http://"+r.Host+"/"+data.GeneratedURL, w)
		} else {
			db.Create(&data)
			utils.FormatResponse("created", "http://"+r.Host+"/"+data.GeneratedURL, w)
		}
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	data := models.StoreURL{}
	response := db.Where(&models.StoreURL{GeneratedURL: path}).First(&data)

	if response.RowsAffected > 0 {
		http.Redirect(w, r, data.OriginalURL, http.StatusSeeOther)
	} else {
		utils.FormatResponse("false", "No result Found", w)
	}
}
