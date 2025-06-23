package api

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Items struct {
	Items []Item `json:"items"`
}
