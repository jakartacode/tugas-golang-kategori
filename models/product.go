package models

type ProductWrite struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	CategoryID int `json:"category_id"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Category string `json:"category"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}