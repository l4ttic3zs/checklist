package main

import (
	"catalog/api"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (a *App) CreateItemType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "Name field is mandatory", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "File missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), handler.Filename)

	uploadPath := "/app/data/images"
	dstPath := filepath.Join(uploadPath, fileName)

	os.MkdirAll(uploadPath, os.ModePerm)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Could not save the file to the server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error while saving file", http.StatusInternalServerError)
		return
	}

	log.Printf("[POST] - [CATALOG] - Saved file to: %s", dstPath)

	newItem := api.ItemType{
		Name:      name,
		ImagePath: "/app/data/images/" + fileName,
	}

	if err := a.DB.Create(&newItem).Error; err != nil {
		http.Error(w, "DB error while saving", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}
