package main

import (
	"catalog/api"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func (a *App) UpdateItemTypeByName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// example /update?name=Michael
	targetName := r.URL.Query().Get("name")
	if targetName == "" {
		http.Error(w, "Missing name parameter", http.StatusBadRequest)
		return
	}

	var item api.ItemType
	if err := a.DB.Where("name = ?", targetName).First(&item).Error; err != nil {
		http.Error(w, "Item not found with this name", http.StatusNotFound)
		return
	}

	r.ParseMultipartForm(10 << 20)
	log.Printf("URL Query Name: %s", r.URL.Query().Get("name"))
	log.Printf("Form Value Name: %s", r.FormValue("new_name"))

	for key, values := range r.MultipartForm.Value {
		log.Printf("Érkező mező: %s = %s", key, values[0])
	}

	if newName := r.FormValue("new_name"); newName != "" {
		item.Name = newName
	}

	file, handler, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		if item.ImagePath != "" {
			os.Remove(item.ImagePath)
		}
		os.MkdirAll("./uploads", os.ModePerm)
		newPath := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), handler.Filename)
		dst, _ := os.Create(newPath)
		defer dst.Close()
		io.Copy(dst, file)
		item.ImagePath = newPath
	}

	if err := a.DB.Save(&item).Error; err != nil {
		http.Error(w, "Conflict: Name already exists or database error", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
