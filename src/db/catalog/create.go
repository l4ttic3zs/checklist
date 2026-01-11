package main

import (
	"catalog/api"
	"encoding/json"
	"log"
	"net/http"
)

func (a *App) CreateItemType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Name      string `json:"name"`
		ImagePath string `json:"imagePath"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[POST] - [CATALOG] - Creating item with name: %s", input.Name)

	newItem := api.ItemType{
		Name:      input.Name,
		ImagePath: input.ImagePath,
	}
	a.DB.Create(&newItem)
	json.NewEncoder(w).Encode(newItem)
}
