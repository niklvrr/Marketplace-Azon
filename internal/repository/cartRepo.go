package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	getCartByUserIdQuery = `
		SELECT id, created_at
		FROM carts
		WHERE user_id = $1;`

	addItemQuery = `
		INSERT INTO cart_items(cart_id, product_id, quantity, created_at)
		VALUES($1, $2, $3, $4)
		RETURNING id;`

	removeItemQuery = `
		DELETE FROM cart_items 
		WHERE cart_id = $1 AND product_id = $2;`

	clearCartQuery = `
		DELETE FROM cart_items
		WHERE cart_id = $1;`
)

var (
	cartNotFoundError     = errors.New("cart not found")
	cartItemNotFoundError = errors.New("cart item not found")
	addItemError          = errors.New("add item error")
	removeItemError       = errors.New("remove item error")
	clearCartError        = errors.New("clear cart error")
)

type CartRepo struct {
	db *pgxpool.Pool
}

func NewCartRepo(db *pgxpool.Pool) *CartRepo {
	return &CartRepo{db: db}
}

func (r *CartRepo) GetCartByUserId(ctx context.Context, userId int64) (*model.Cart, error) {
	cart := new(model.Cart)
	err := r.db.QueryRow(ctx, getCartByUserIdQuery, userId).Scan(&cart.Id, &cart.CreatedAt)
	if err != nil {
		return nil, cartNotFoundError
	}

	return cart, nil
}

func (r *CartRepo) AddItem(ctx context.Context, cartId, productId int64, quantity int) (int64, error) {
	var itemId int64

	err := r.db.QueryRow(
		ctx, addItemQuery,
		cartId, productId, quantity).Scan(&itemId)
	if err != nil {
		return 0, addItemError
	}

	return itemId, nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, cartId, productId int64) error {
	cmdTag, err := r.db.Exec(ctx, removeItemQuery, cartId, productId)
	if err != nil {
		return removeItemError
	}

	if cmdTag.RowsAffected() == 0 {
		return cartItemNotFoundError
	}

	return nil
}

func (r *CartRepo) ClearCart(ctx context.Context, cartId int64) error {
	cmdTag, err := r.db.Exec(ctx, clearCartQuery, cartId)
	if err != nil {
		return clearCartError
	}

	if cmdTag.RowsAffected() == 0 {
		return cartNotFoundError
	}

	return nil
}
