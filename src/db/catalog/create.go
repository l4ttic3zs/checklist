package main

import (
	"catalog/api"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func (a *App) CreateItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	name := r.FormValue("name")

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "You must upload an image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	os.MkdirAll("./uploads", os.ModePerm)
	filePath := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), handler.Filename)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error during image save", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	newItem := api.ItemType{
		Name:      name,
		ImagePath: filePath,
	}
	a.DB.Create(&newItem)
	json.NewEncoder(w).Encode(newItem)
}
