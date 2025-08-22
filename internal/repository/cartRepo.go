package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	getCartByUserIdQuery = `
		SELECT id, created_at
		FROM carts
		WHERE user_id = $1;`

	getCartItemsByCartIdQuery = `
		SELECT id, product_id, quantity
		FROM cart_items 
		WHERE cart_id = $1;`

	addItemQuery = `
		INSERT INTO cart_items(cart_id, product_id, quantity, created_at)
		VALUES($1, $2, $3, $4)
		RETURNING id;`

	removeItemQuery = `
		DELETE FROM cart_items 
		WHERE cart_id = $1 AND id = $2;`

	clearCartQuery = `
		DELETE FROM cart_items
		WHERE cart_id = $1;`
)

var (
	cartNotFoundError         = errors.New("cart not found")
	cartItemNotFoundError     = errors.New("cart item not found")
	addItemError              = errors.New("add item error")
	getCartItemsByCartIdError = errors.New("cart item by cart id error")
	removeItemError           = errors.New("remove item error")
	clearCartError            = errors.New("clear cart error")
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
		return nil, fmt.Errorf("%w: %w", cartNotFoundError, err)
	}

	return cart, nil
}

func (r *CartRepo) GetCartItemsByCartId(ctx context.Context, cartId int64) (*[]model.CartItem, error) {
	rows, err := r.db.Query(ctx, getCartItemsByCartIdQuery, cartId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", getCartItemsByCartIdError, err)
	}
	defer rows.Close()

	var cartItems []model.CartItem
	for rows.Next() {
		var cartItem model.CartItem
		err = rows.Scan(
			&cartItem.Id,
			&cartItem.ProductId,
			&cartItem.Quantity,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %w", getCartItemsByCartIdQuery, err)
		}

		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w(%w): %w", getCartItemsByCartIdError, rowsIterationError, err)
	}

	return &cartItems, nil
}

func (r *CartRepo) AddItem(ctx context.Context, cartId, productId int64, quantity int) (int64, error) {
	var itemId int64

	err := r.db.QueryRow(
		ctx, addItemQuery,
		cartId, productId, quantity).Scan(&itemId)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", addItemError, err)
	}

	return itemId, nil
}

func (r *CartRepo) RemoveItem(ctx context.Context, cartId, id int64) error {
	cmdTag, err := r.db.Exec(ctx, removeItemQuery, cartId, id)
	if err != nil {
		return fmt.Errorf("%w: %w", removeItemError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", removeItemError, cartItemNotFoundError)
	}

	return nil
}

func (r *CartRepo) ClearCart(ctx context.Context, cartId int64) error {
	cmdTag, err := r.db.Exec(ctx, clearCartQuery, cartId)
	if err != nil {
		return fmt.Errorf("%w: %w", clearCartError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", clearCartError, cartNotFoundError)
	}

	return nil
}
