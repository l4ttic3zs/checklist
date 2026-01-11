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

	var item api.ShoppingList
	if err := a.DB.Where("name = ?", targetName).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
	}
	if err := a.DB.Delete(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
