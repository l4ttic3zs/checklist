package main

import (
	"encoding/json"
	"fmt"
	"inventory/api"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (a *App) StartListening() {
	ch, err := a.RMQ.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %e", err)
	}

	msgs, err := ch.Consume("purchase_queue", "", true, false, false, false, nil)
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
			if err := a.updateInventory(data.Name, data.Count); err != nil {
				log.Printf("ERROR updating inventory: %v", err)
			}
		}
	}()
}

func (a *App) updateInventory(name string, count int) error {
	var itemType api.ItemType
	if err := a.DB.Where("name = ?", name).First(&itemType).Error; err != nil {
		return fmt.Errorf("item type '%s' not found", name)
	}

	err := a.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "item_type_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"count": gorm.Expr("items.count + ?", count),
		}),
	}).Create(&api.Item{
		ItemTypeID: itemType.ID,
		Count:      count,
	}).Error

	if err != nil {
		return fmt.Errorf("upsert failed: %v", err)
	}

	log.Printf("Inventory updated: %s (+%d)", name, count)
	return nil
}
