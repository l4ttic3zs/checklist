package app

import (
	"checklist/pkg/api"
	"database/sql"
)

type App struct {
	db    *sql.DB
	Items api.Items
}

func NewApp() (*App, error) {
	return &App{Items: api.Items{}}, nil
}
