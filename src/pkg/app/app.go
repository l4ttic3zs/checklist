package app

import (
	"checklist/pkg/api"
	"log"
)

func (app *App) GetItems() (api.Items, error) {

	return api.Items{}, nil
}

func (app *App) AddNewItem(item api.Item) (int, error) {
	if err := app.db.Ping(); err != nil {
		log.Println("Error pinging db")
	}

	insertItemSQL := `INSERT INTO items (name, quantity) VALUES($1, $2) RETURNING id;`
	var lastInsertId int
	err := app.db.QueryRow(insertItemSQL, item.Name, item.Quantity).Scan(lastInsertId)
	if err != nil {
		log.Println("Error inserting new item")
	}

	return lastInsertId, nil
}
