package main

import (
	"checklist/pkg/app"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/api", app.TestApi)

	err := http.ListenAndServe(":8443", router)
	if err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
