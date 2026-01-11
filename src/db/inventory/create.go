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
		ItemTypeName string `json:"item_type"`
		Count        int    `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var itemType api.ItemType
	if err := a.DB.Where("name = ?", input.ItemTypeName).First(&itemType).Error; err != nil {
		http.Error(w, "Item type not found in catalog. Create it there first!", http.StatusNotFound)
		return
	}

	var existingItem api.Item
	err := a.DB.Where("item_type_id = ?", itemType.ID).First(&existingItem).Error
	if err == nil {
		http.Error(w, "This item already exist from this type. Use PUT to modify it.", http.StatusConflict)
		return
	}

	newItem := api.Item{
		Count: input.Count,
	}

	if err := a.DB.Create(&newItem).Error; err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	a.DB.Preload("ItemType").First(&newItem, newItem.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}
