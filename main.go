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
		log.Fatal("Ecannot load environment file")
	}
}

func connectDB() (*gorm.DB, error) {
	DBUsername := os.Getenv("DB_USERNAME")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBName := os.Getenv("DB_NAME")
	DBHost := os.Getenv("DB_HOST")
	conn := fmt.Sprintf("%s:%s@(%s)/%s", DBUsername, DBPassword, DBHost, DBName)

	db, err := gorm.Open("mysql", conn)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}
	log.Println("established connection to database")
	return db, nil

}

func migrateSchema() error {
	dbConn, err := connectDB()
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	dbConn.AutoMigrate(&StoreURL{}) //refactor to init or migrate function called in main, create a function to handle this and return error.
	handleRequests(dbConn)
	return nil
}

func main() {
	err := migrateSchema()
	if err != nil {
		fmt.Errorf("cannot migrate schema: %v", err)
	}
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
