package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shoppinglist/api"

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

	err = db.AutoMigrate(&api.ItemType{}, &api.ShoppingList{})
	if err != nil {
		log.Printf("Migration error: %v", err)
	}

	app := &App{DB: db}

	log.Println("Starting server on port 80...")
	http.HandleFunc("/shoppinglist", app.GetShoppingListItems)
	http.HandleFunc("/shoppinglist/item", app.HandleItem)

	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func (a *App) HandleItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.GetShoppingListItems(w, r)
	case http.MethodPost:
		a.CreateShoppingListItem(w, r)
	case http.MethodPut:
		a.UpdateShoppingListItemByName(w, r)
	case http.MethodDelete:
		a.DeleteShoppingListItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
