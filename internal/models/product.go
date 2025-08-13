package models

import "time"

type Product struct {
	Id          int64     `json:"id" db:"id"`
	SellerId    int64     `json:"seller_id" db:"seller_id"`
	CategoryId  int64     `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Stock       int       `json:"stock" db:"stock"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
