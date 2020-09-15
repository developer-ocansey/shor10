package routes

import (
	"net/http"

	"urlshortner/controllers"

	"github.com/gorilla/mux"
)

// HandleRequests ..
func HandleRequests() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/shorten-url", controllers.ShortenURL).Methods("POST") // POST long URL and get short path
	router.HandleFunc("/{path}", controllers.Redirect).Methods("GET")         // Pass value to URL and redirect to original path stored

	router.HandleFunc("/", controllers.Healthz).Methods("GET")
	router.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	router.HandleFunc("/login", controllers.Login).Methods("POST")

	// Auth route
	// subrouter := router.PathPrefix("/auth").Subrouter()
	// subrouter.Use(auth.JwtVerify)
	// subrouter.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	// subrouter.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	// subrouter.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	// subrouter.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")

	return router
}

// CommonMiddleware ..
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
