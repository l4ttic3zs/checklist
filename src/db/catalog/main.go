package main

import (
	"catalog/api"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	DB *gorm.DB
}

func main() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	err = db.AutoMigrate(&api.Item{})
	if err != nil {
		log.Printf("Migration error: %v", err)
	}

	app := &App{DB: db}

	log.Println("Starting server on port 80...")
	http.HandleFunc("/list", app.GetItems)
	http.HandleFunc("/upload", app.CreateItem)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func (a *App) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var items []api.Item
	result := a.DB.Find(&items)
	if result.Error != nil {
		log.Printf("Error during query: %v", result.Error)
		http.Error(w, fmt.Sprintf("Database error: %v", result.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (a *App) CreateItem(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	name := r.FormValue("name")

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "You must upload an image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("./uploads", os.ModePerm)
	filePath := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error during image save", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	newItem := api.Item{
		Name:      name,
		ImagePath: filePath,
	}
	a.DB.Create(&newItem)

	json.NewEncoder(w).Encode(newItem)
}
