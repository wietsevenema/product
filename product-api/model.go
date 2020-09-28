package model

type Product struct {
	ProductID   string `json:"sku"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Price       int    `json:"price"`
	Upc         string `json:"upc"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	URL         string `json:"url"`
	Image       string `json:"image"`
	Category    string `json:"category"`
}
