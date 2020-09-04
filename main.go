package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

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
	URL string `json:"url"`
}

// Handler ...
type Handler struct {
	db *gorm.DB
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	dbConString, err := goDotEnvVariable("DBCONNECTINGSTRING")
	if err != nil {
		fmt.Errorf("cannot call goDotEnvVariable function: ", err)
	}

	db, err := gorm.Open("mysql", dbConString)

	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&StoreURL{}) //refactor to init or migrate function called in main
	handleRequests(db)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	formatResponse("true", "Success is not final; failure is not fatal: It is the courage to continue that counts.", w)
}

func (h *Handler) shortenURL(w http.ResponseWriter, r *http.Request) {
	var bodyURL URL

	if err := json.NewDecoder(r.Body).Decode(&bodyURL); err != nil {
		log.Fatal(err)
	}

	url := bodyURL.URL
	if len(url) < 25 {
		formatResponse("invalid", "URL must be more than 24 characters", w)
	} else {
		hash := sha1.New()
		hash.Write([]byte(url))
		hashedURL := string(hash.Sum(nil))
		sliceURL := hashedURL[0:3]
		shortURL := fmt.Sprintf("%x", sliceURL)

		data := StoreURL{GeneratedURL: shortURL, OriginalURL: url}
		response := h.db.Where(&StoreURL{OriginalURL: url}).Find(&data)

		if response.RowsAffected > 0 {
			formatResponse("exist", "http://"+r.Host+"/"+data.GeneratedURL, w)
		} else {
			h.db.Create(&data)
			formatResponse("created", "http://"+r.Host+"/"+data.GeneratedURL, w)
		}
	}
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]
	data := StoreURL{}
	response := h.db.Where(&StoreURL{GeneratedURL: path}).First(&data)

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

func handleRequests(db *gorm.DB) {
	handler := &Handler{db: db}
	log.Println("Starting development server at http://localhost:8020/")
	log.Println("Quit the server with CONTROL-C.")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", healthz)
	router.HandleFunc("/shorten-url", handler.shortenURL) // POST long URL and get short path
	router.HandleFunc("/{path}", handler.redirect)        // Pass value to URL and redirect to original path stored
	log.Fatal(http.ListenAndServe(":8020", router))
}

func goDotEnvVariable(key string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", fmt.Errorf("Error loading .env file")
	}
	return os.Getenv(key), nil
}
