package main

import (
	"net/http"
	"shoppinglist/api"
)

func (a *App) DeleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	targetName := r.URL.Query().Get("name")
	if targetName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	var item api.ShoppingList
	err := a.DB.Joins("ItemType").
		Where("\"ItemType\".name = ?", targetName).
		First(&item).Error
	if err != nil {
		http.Error(w, "Item not found with the given type name", http.StatusNotFound)
		return
	}

	if err := a.DB.Delete(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
