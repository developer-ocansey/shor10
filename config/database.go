package config

import (
	"fmt"
	"log"
	"os"

	"urlshortner/models"

	"github.com/jinzhu/gorm"
)

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
	// close db when not in use
	defer db.Close()
	log.Println("established connection to database")
	return db, nil

}

// InitDB ...
func InitDB() (*gorm.DB, error) {
	dbConn, err := connectDB()
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	dbConn.AutoMigrate(&models.StoreURL{}, &models.User{}) //refactor to init or migrate function called in main, create a function to handle this and return error.
	return dbConn, nil
}
