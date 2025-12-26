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

	log.Printf("[DELETE] - [CATALOG] - Deleting item by name: %s", targetName)

	var itemtype api.ItemType
	if err := a.DB.Where("name = ?", targetName).First(&itemtype).Error; err != nil {
		http.Error(w, "Itemtype not found", http.StatusNotFound)
		return
	}

	if itemtype.ImagePath != "" {
		err := os.Remove(itemtype.ImagePath)
		if err != nil {
			log.Printf("Could not delete file: %v", err)
		}
	}

	if err := a.DB.Delete(&itemtype).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
