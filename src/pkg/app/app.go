package app

import "checklist/pkg/api"

func (app *App) GetItems() (api.Items, error) {

	return api.Items{}, nil
}

func (app *App) AddNewItem(item api.Item) error {
	// call db add
	return nil
}
