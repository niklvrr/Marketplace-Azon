package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createOrderQuery = `
		INSERT INTO orders (user_id, status, total, create_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	createOrderItemQuery = `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	getOrdersByUserIdQuery = `
		SELECT order_id, status, total, create_at
		FROM orders
		WHERE user_id = $1`

	getOrderByIdQuery = `
		SELECT user_id, status, total, create_at
		FROM orders
		WHERE order_id = $1`

	getOrderItemsByOrderIdQuery = `
		SELECT id, product_id, quantity, price
		FROM order_items
		WHERE order_id = $1`

	deleteOrderByIdQuery = `DELETE FROM orders WHERE id = $1`
)

var (
	createOrderError            = errors.New("error creating order")
	createOrderItemError        = errors.New("error creating orderItem")
	orderNotFound               = errors.New("order not found")
	getOrdersByUserIdError      = errors.New("error getting orders by user id")
	getOrderByIdError           = errors.New("error getting order by id")
	getOrderItemsByOrderIdError = errors.New("error getting order items by order id")
	deleteOrderByIdError        = errors.New("error deleting order by id")
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error) {
	var total float64
	for _, orderItem := range *items {
		total += orderItem.Price
	}

	var orderId int64
	for _, item := range *items {
		err := r.db.QueryRow(
			ctx, createOrderQuery,
			userId,
			"pending",
			total,
			time.Now()).Scan(&orderId)

		if err != nil {
			return 0, fmt.Errorf("%w: %w", createOrderError, err)
		}

		err = r.db.QueryRow(
			ctx, createOrderItemQuery,
			orderId,
			item.ProductId,
			item.Quantity,
			item.Price,
		).Scan(&item.Id)

		if err != nil {
			return 0, fmt.Errorf("%w: %w", createOrderItemError, err)
		}
	}

	return orderId, nil
}

func (r *OrderRepo) GetOrdersByUserId(ctx context.Context, userId int64) (*[]model.Order, error) {
	rows, err := r.db.Query(ctx, getOrdersByUserIdQuery, userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", getOrdersByUserIdError, err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order

		err := rows.Scan(
			&order.Id,
			&order.Status,
			&order.Total,
			&order.CreateAt,
		)
		order.UserId = userId

		if err != nil {
			return nil, fmt.Errorf("%w: %w", getOrdersByUserIdError, err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w(%w): %w", getOrdersByUserIdError, rowsIterationError, err)
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrderById(ctx context.Context, orderId int64) (*model.Order, error) {
	order := new(model.Order)
	order.Id = orderId
	err := r.db.QueryRow(ctx, getOrderByIdQuery, orderId).Scan(&order.UserId, &order.Status, &order.Total, &order.CreateAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", getOrderByIdError, err)
	}

	return order, nil
}

func (r *OrderRepo) GetOrderItemsByOrderId(ctx context.Context, orderId int64) (*[]model.OrderItem, error) {
	rows, err := r.db.Query(ctx, getOrderItemsByOrderIdQuery, orderId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", getOrderItemsByOrderIdError, err)
	}
	defer rows.Close()

	var orderItems []model.OrderItem
	for rows.Next() {
		var orderItem model.OrderItem
		err := rows.Scan(
			&orderItem.Id,
			&orderItem.ProductId,
			&orderItem.Quantity,
			&orderItem.Price,
		)

		if err != nil {
			return nil, fmt.Errorf("%w: %w", getOrderItemsByOrderIdError, err)
		}

		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w(%w): %w", getOrderItemsByOrderIdError, rowsIterationError, err)
	}

	return &orderItems, nil
}

func (r *OrderRepo) DeleteOrderById(ctx context.Context, orderId int64) error {
	cmdTag, err := r.db.Exec(ctx, deleteOrderByIdQuery, orderId)
	if err != nil {
		return fmt.Errorf("%w: %w", deleteOrderByIdError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", deleteOrderByIdError, orderNotFound)
	}

	return nil
}
