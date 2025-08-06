package models

type CartItem struct {
	Id        int64 `json:"id" db:"id"`
	CartId    int64 `json:"cart_id" db:"cart_id"`
	ProductId int64 `json:"product_id" db:"product_id"`
	Quantity  int64 `json:"quantity" db:"quantity"`
}
