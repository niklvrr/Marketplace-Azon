package models

type OrderItem struct {
	Id        int64   `json:"id" db:"id"`
	OrderId   int64   `json:"order_id" db:"order_id"`
	ProductId int64   `json:"product_id" db:"product_id"`
	Quantity  int64   `json:"quantity" db:"quantity"`
	Price     float64 `json:"price" db:"price"`
}
