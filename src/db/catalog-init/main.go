package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name string
}

func main() {
	migrateMode := flag.Bool("migrate", false, "Run migration mod")
	flag.Parse()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error while connecting to DB: %v", err)
	}

	if *migrateMode {
		fmt.Println("Starting migration...")
		err := db.AutoMigrate(&Item{})
		if err != nil {
			log.Fatalf("Migration error: %v", err)
		}
		fmt.Println("Migration successfull.")
		os.Exit(0)
	}
}
