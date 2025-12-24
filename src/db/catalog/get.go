package main

import (
	"catalog/api"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (a *App) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println("Getting items")

	var items []api.ItemType
	result := a.DB.Find(&items)
	if result.Error != nil {
		log.Printf("Error during query: %v", result.Error)
		http.Error(w, fmt.Sprintf("Database error: %v", result.Error), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (a *App) GetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// example /update?name=Michael
	targetName := r.URL.Query().Get("name")
	if targetName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	var item api.ItemType
	if err := a.DB.Where("name = ?", targetName).First(&item).Error; err != nil {
		http.Error(w, "Item not found with this name", http.StatusNotFound)
		return
	}

	json, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "Error during item marshalling", http.StatusInternalServerError)
		return
	}
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}
