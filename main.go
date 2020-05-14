package main

import (
	b64 "encoding/base64"
	"encoding/json"
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
	ID          int    `json:"id"`
	FakeURL     string `json:"fakeURL"` // generatedURL...
	OriginalURL string `json:"originalURL"`
}

// Response ...
type Response struct {
	status     bool
	shortenURL string
	error      string
}

// Value ...
type Value struct {
	id          int
	fakeURL     string
	originalURL string
}

// Query ...
type Query struct {
	Value        Value
	Error        string
	RowsAffected int
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

	err := r.ParseForm()
	if err != nil {
		formatResponse("false", "Kindly Enter a value for url", w)
	}
	url := r.FormValue("url")
	// u, err := url.ParseRequestURI(url)
	// if err != nil {
	// 	formatResponse("false", "Invalid URL", w)
	// }

	bs64 := b64.StdEncoding.EncodeToString([]byte(url))
	shortURL := bs64[len(bs64)-6 : len(bs64)]
	data := StoreURL{FakeURL: shortURL, OriginalURL: url}
	response := db.Where(&StoreURL{OriginalURL: url}).Find(&data)

	if response.RowsAffected > 0 {
		formatResponse("exist", r.Host+"/"+data.FakeURL, w)
	} else {
		db.Create(&data)
		formatResponse("created", r.Host+"/"+data.FakeURL, w)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	path := vars["path"]
	data := StoreURL{}
	response := db.Where(&StoreURL{FakeURL: path}).First(&data)

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
}
