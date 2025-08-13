package models

// Product models
type GetProductsRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type CreateProductRequest struct {
	CategoryId  string  `json:"category_id" binding:"required"`
	Name        string  `json:"name" binding:"required, min=2, max=100"`
	Description string  `json:"description" binding:"omitempty,max=5000"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"required, min=0"`
}

type UpdateProductRequest struct {
	CategoryId  *string  `json:"category_id" binding:"required"`
	Name        *string  `json:"name" binding:"required, min=2, max=100"`
	Description *string  `json:"description" binding:"omitempty,max=5000"`
	Price       *float64 `json:"price" binding:"required,gt=0"`
	Stock       *int     `json:"stock" binding:"required, min=0"`
}

type DeleteProductRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type SearchProductsRequest struct {
	Text       string   `form:"text" binding:"omitempty,min=1,max=100"`
	CategoryId *string  `form:"category_id" binding:"omitempty"`
	Min        *float64 `form:"min" binding:"omitempty,gt=0"`
	Max        *float64 `form:"max" binding:"omitempty,gt=0"`
	Offset     int64    `form:"offset" binding:"omitempty,gt=0"`
	Limit      int64    `form:"limit" binding:"omitempty,gt=0"`
}

// User models
type SighUpRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Admin models
type GetUserByIdRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type ApproveProductRequest struct {
	ProductId int64 `json:"product_id" binding:"required"`
	IsApprove bool  `json:"is_approve" binding:"required"`
}

// Cart models
type AddItemRequest struct {
	CartId    int64 `json:"cart_id" binding:"required"`
	ProductId int64 `json:"product_id" binding:"required"`
	Quantity  int64 `json:"quantity" binding:"required"`
}

type RemoveItemRequest struct {
	CartId    int64 `json:"cart_id" binding:"required"`
	ProductId int64 `json:"product_id" binding:"required"`
}

type ClearCartRequest struct {
	CartId int64 `json:"cart_id" binding:"required"`
}

type GetCartByUserIdRequest struct {
	UserId int64 `json:"user_id" binding:"required"`
}

// Order models
type CreateOrderRequest struct {
	Status string  `json:"status" binding:"required"`
	Total  float64 `json:"total" binding:"required"`
}

type GetOrderByIdRequest struct {
	Id int64 `json:"id" binding:"required"`
}

// Category models
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"omitempty,max=5000"`
}

type GetCategoryByIdRequest struct {
	Id int64 `json:"id" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"required,min=2,max=100"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
}

type DeleteCategoryRequest struct {
	Id int64 `json:"id" binding:"required"`
}
