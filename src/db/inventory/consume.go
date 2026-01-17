package main

import (
	"encoding/json"
	"fmt"
	"inventory/api"
	"log"
)

func (a *App) StartListening() {
	ch, err := a.RMQ.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %e", err)
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
        log.Fatalf("failed to declare a queue: %v", err)
    }

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%e", err)
	}

	go func() {
		for d := range msgs {
			var data struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}
			json.Unmarshal(d.Body, &data)

			log.Printf("Purchase arrived: %d piece", data.Count)
			a.updateInventory(data.Name, data.Count)
		}
	}()
}

func (a *App) updateInventory(name string, count int) error {
	var itemType api.ItemType
	if err := a.DB.Where("name = ?", name).First(&itemType).Error; err != nil {
		return fmt.Errorf("Item type not found with this name")
	}

	var item api.Item
	if err := a.DB.Where("item_type_id = ?", itemType.ID).First(&item).Error; err != nil {
		return fmt.Errorf("Item not found in inventory")
	}

	item.Count += count
	if err := a.DB.Save(&item).Error; err != nil {
		return fmt.Errorf("Could not save")
	}
	return nil
}
