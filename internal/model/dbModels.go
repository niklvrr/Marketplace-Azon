package model

import "time"

type Cart struct {
	Id        int64     `json:"id" db:"id"`
	UserId    int64     `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CartItem struct {
	Id        int64 `json:"id" db:"id"`
	CartId    int64 `json:"cart_id" db:"cart_id"`
	ProductId int64 `json:"product_id" db:"product_id"`
	Quantity  int64 `json:"quantity" db:"quantity"`
}

type Category struct {
	Id          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

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

type Order struct {
	Id       int64     `json:"id" db:"id"`
	UserId   int64     `json:"user_id" db:"user_id"`
	Status   string    `json:"status" db:"status"`
	Total    float64   `json:"total" db:"total"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}

type OrderItem struct {
	Id        int64   `json:"id" db:"id"`
	OrderId   int64   `json:"order_id" db:"order_id"`
	ProductId int64   `json:"product_id" db:"product_id"`
	Quantity  int64   `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}

type User struct {
	Id       int64     `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
	Role     string    `json:"role" db:"role"`
	IsActive bool      `json:"is_active" db:"is_active"`
	CreateAt time.Time `json:"create_at" db:"create_at"`
}
