package main

import (
	"catalog/api"
	"encoding/json"
	"log"
	"net/http"
)

func (a *App) CreateItemType(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20)
	name := r.FormValue("name")
	filePath := r.FormValue("path")

	log.Printf("[POST] - [CATALOG] - Creating item with name: %s", name)

	newItem := api.ItemType{
		Name:      name,
		ImagePath: filePath,
	}
	a.DB.Create(&newItem)
	json.NewEncoder(w).Encode(newItem)
}
