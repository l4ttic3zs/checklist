package main

import (
	"encoding/json"
	"fmt"
	"inventory/api"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (a *App) StartListening() {
	conn, _ := amqp.Dial("amqp://guest:guest@checklist-purchase-queue:5672/")
	ch, _ := conn.Channel()

	msgs, _ := ch.Consume("purchase_queue", "", true, false, false, false, nil)

	go func() {
		for d := range msgs {
			var data struct {
				Name  string `json:"name"`
				Count int    `json:"count"`
			}
			json.Unmarshal(d.Body, &data)

			log.Printf("Vásárlás érkezett: %d darab", data.Count)
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
