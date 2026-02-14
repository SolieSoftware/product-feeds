package models

import "time"

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	ImagePaths  []string  `json:"image_paths"`
	SourceURL   string    `json:"source_url"`
	CompanyName string    `json:"company_name"`
	ScrapedAt   time.Time `json:"scraped_at"`
}
