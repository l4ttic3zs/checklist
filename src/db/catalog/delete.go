package main

import (
	"catalog/api"
	"log"
	"net/http"
	"os"
)

func (a *App) DeleteItemTypeByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	targetName := r.URL.Query().Get("name")
	if targetName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	var item api.ItemType
	if err := a.DB.Where("name = ?", targetName).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	if item.ImagePath != "" {
		err := os.Remove(item.ImagePath)
		if err != nil {
			log.Printf("Could not delete file: %v", err)
		}
	}

	a.DB.Delete(&item)
	w.WriteHeader(http.StatusOK)
}
