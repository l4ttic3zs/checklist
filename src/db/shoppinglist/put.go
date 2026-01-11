package main

import (
	"encoding/json"
	"net/http"
	"shoppinglist/api"
)

func (a *App) UpdateShoppingListItemByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	var item api.ShoppingList
	if err := a.DB.Where("item_type_id = ?", itemType.ID).First(&item).Error; err != nil {
		http.Error(w, "Item not found in inventory", http.StatusNotFound)
		return
	}

	var input struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	item.Count = input.Count
	if err := a.DB.Save(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item.ItemType = itemType
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}
