package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func (app *App) ConnectDB() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	app.db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error during connection: %v", err)
	}

	err = app.db.Ping()
	if err != nil {
		log.Fatalf("Error: db doesn't respong: %v", err)
	}
	log.Println("Connected to DB")

	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		quantity SMALLINT NOT NULL
	);`
	_, err = app.db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error making db table: %v", err)
	}
}

func (app *App) CloseDB() {
	app.db.Close()
}
