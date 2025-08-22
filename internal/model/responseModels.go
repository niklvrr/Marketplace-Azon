package model

type ProductResponse struct {
	Id          int64   `json:"id"`
	SellerId    int64   `json:"seller_id"`
	CategoryId  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type UserResponse struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CartResponse struct {
	Id int64 `json:"id"`
}

type CartItemResponse struct {
	Id        int64 `json:"id"`
	CartId    int64 `json:"cart_id"`
	ProductId int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CategoryResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrderResponse struct {
	Id     int64   `json:"id"`
	UserId int64   `json:"user_id"`
	Status string  `json:"status"`
	Total  float64 `json:"total"`
}

type OrderItemResponse struct {
	Id        int64   `json:"id"`
	OrderId   int64   `json:"order_id"`
	ProductId int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Price     float64 `json:"price"`
}
