package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// StoreURL ...
type StoreURL struct {
	ID           int    `json:"id"`
	GeneratedURL string `json:"generatedURL"`
	OriginalURL  string `json:"originalURL"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

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
	url string
}

func handleRequests() {
	log.Println("Starting development server at http://localhost:8020/")
	log.Println("Quit the server with CONTROL-C.")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", healthz)
	myRouter.HandleFunc("/shorten-url", shortenURL) //POST long URL and get short path
	myRouter.HandleFunc("/{path}", redirect)        // pass value to URL and redirect to origin al path stored
	log.Fatal(http.ListenAndServe(":8020", myRouter))
}

func main() {
	db, err = gorm.Open("mysql", "b9d1984e35796d:3d04ee89@(eu-cdbr-west-03.cleardb.net)/heroku_f12d52ab1d22e5f")

	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&StoreURL{})
	handleRequests()
}

func healthz(w http.ResponseWriter, r *http.Request) {
	formatResponse("true", "Success is not final; failure is not fatal: It is the courage to continue that counts.", w)
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	var bodyURL URL

	if err := json.NewDecoder(r.Body).Decode(&bodyURL); err != nil {
		log.Fatal(err)
	}

	url := bodyURL.url

	h := sha1.New()
	h.Write([]byte(url))
	hashedURL := string(h.Sum(nil))
	sliceURL := hashedURL[len(hashedURL)-3 : len(hashedURL)]
	shortURL := fmt.Sprintf("%x", sliceURL)

	data := StoreURL{GeneratedURL: shortURL, OriginalURL: url}
	response := db.Where(&StoreURL{OriginalURL: url}).Find(&data)

	if response.RowsAffected > 0 {
		formatResponse("exist", r.Host+"/"+data.GeneratedURL, w)
	} else {
		db.Create(&data)
		formatResponse("created", r.Host+"/"+data.GeneratedURL, w)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	data := StoreURL{}
	response := db.Where(&StoreURL{GeneratedURL: path}).First(&data)

	if response.RowsAffected > 0 {
		http.Redirect(w, r, data.OriginalURL, http.StatusSeeOther)
	} else {
		formatResponse("false", "No result Found", w)
	}
}

func formatResponse(status string, message string, w http.ResponseWriter) {
	res := make(map[string]string)
	res["status"] = status
	res["message"] = message
	json.NewEncoder(w).Encode(res)
	return
}
