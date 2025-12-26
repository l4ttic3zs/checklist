package main

import (
	"encoding/json"
	"inventory/api"
	"net/http"
)

func (a *App) CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		ItemTypeID uint `json:"item_type_id"`
		Count      int  `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newItem := api.Item{
		ItemTypeID: input.ItemTypeID,
		Count:      input.Count,
	}

	if err := a.DB.Create(&newItem).Error; err != nil {
		http.Error(w, "Could not create item (check if item_type_id exists)", http.StatusInternalServerError)
		return
	}

	// Preload-oljuk a választ, hogy a kliens lássa a típus nevét is
	a.DB.Preload("ItemType").First(&newItem, newItem.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}
