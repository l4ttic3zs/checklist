package main

import (
	"encoding/json"
	"inventory/api"
	"net/http"
)

func (a *App) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var items []api.Item
	result := a.DB.Preload("ItemType").Find(&items)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (a *App) GetItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	var itemType api.ItemType
	if err := a.DB.Where("name = ?", name).First(&itemType).Error; err != nil {
		http.Error(w, "Item type not found with this name", http.StatusNotFound)
		return
	}

	var item api.Item
	if err := a.DB.Where("item_type_id = ?", itemType.ID).First(&item).Error; err != nil {
		http.Error(w, "Item not found in inventory", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
