package handlers

import (
	"checklist/pkg/api"
	"encoding/json"
	"net/http"
)

func TestApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func AddNewItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body api.Item
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// err = app.AddNewItem(body)
	// if err != nil {
	// 	log.Fatalf("err")
	// }
}
