package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"urlshortner/routes"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()

	if e != nil {
		log.Fatal("cannot load environment file")
	}
	fmt.Println(e)

	port := os.Getenv("PORT")

	// Handle routes
	http.Handle("/", routes.HandleRequests())

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
