package main

import (
	"checklist/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	// app, err := application.NewApp()
	// if err != nil {
	// 	log.Fatalf("Error creating app: %v", err)
	// }

	router := http.NewServeMux()

	router.HandleFunc("/api", handlers.TestApi)
	// router.HandleFunc("/api/getitems")
	router.HandleFunc("/api/addnewitem", handlers.AddNewItemHandler)

	log.Println("Server listening on http://:8443")
	err := http.ListenAndServe(":8443", router)
	if err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
