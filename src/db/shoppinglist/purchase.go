package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"shoppinglist/api"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (a *App) ItemPurchased(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if input.Name == "" || input.Count <= 0 {
		http.Error(w, "Missing name or invalid number", http.StatusBadRequest)
		return
	}

	var itemType api.ItemType
	if err := a.DB.Where("name = ?", input.Name).First(&itemType).Error; err != nil {
		http.Error(w, "Item type not found with this name", http.StatusNotFound)
		return
	}

	var item api.ShoppingList
	if err := a.DB.Where("item_type_id = ?", itemType.ID).First(&item).Error; err != nil {
		http.Error(w, "Item not found in inventory", http.StatusNotFound)
		return
	}

	item.Count = 0
	if err := a.DB.Delete(&item).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := a.SendPurchase(input.Count, input.Name); err != nil {
		log.Printf("RabbitMQ error: %v", err)
		http.Error(w, "Error while creating message", http.StatusInternalServerError)
		return
	}

	item.ItemType = itemType
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
	w.WriteHeader(http.StatusAccepted)
}

func (a *App) SendPurchase(count int, name string) error {
	ch, err := a.RMQ.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"purchase_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{
		"name":  name,
		"count": count,
	})

	log.Printf("Publishing: %s", body)
	return ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
