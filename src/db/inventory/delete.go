package main

import (
	"inventory/api"
	"net/http"
)

func (a *App) DeleteItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	targetName := r.URL.Query().Get("name")

	var item api.Item
	if err := a.DB.Where("name = ?", targetName).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
	}
	if err := a.DB.Delete(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
