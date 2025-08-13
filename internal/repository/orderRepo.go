package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	models2 "github.com/niklvrr/myMarketplace/internal/models"
)

var (
	createOrderQuery = `
		INSERT INTO orders (user_id, status, total, create_at)
		VALUES ($1, $2, $3, $4)
		RETURNING order_id`

	getOrdersByUserIdQuery = `
		SELECT order_id, status, total, create_at
		FROM orders
		WHERE user_id = $1`

	getOrderByIdQuery = `
		SELECT user_id, status, total, create_at
		FROM orders
		WHERE order_id = $1`
)

var (
	createOrderError       = errors.New("error creating order")
	orderByUserIdNotFound  = errors.New("order not found")
	getOrdersByUserIdError = errors.New("error getting orders by user id")
	getOrderByIdError      = errors.New("error getting order by id")
)

type OrderRepo struct {
	db *pgxpool.Pool
}

func (r *OrderRepo) CreateOrder(ctx context.Context, o *models2.Order, items []models2.OrderItem) error {
	for _, item := range items {
		err := r.db.QueryRow(
			ctx, createOrderQuery,
			o.UserId,
			o.Status,
			o.Total,
			o.CreateAt).Scan(&item.OrderId)

		if err != nil {
			return createOrderError
		}
	}

	return nil
}

func (r *OrderRepo) GetOrdersByUserId(ctx context.Context, userId int64) (*[]models2.Order, error) {
	rows, err := r.db.Query(ctx, getOrdersByUserIdQuery, userId)
	if err != nil {
		return nil, getOrdersByUserIdError
	}
	defer rows.Close()

	var orders []models2.Order
	for rows.Next() {
		var order models2.Order

		err := rows.Scan(
			&order.Id,
			&order.Status,
			&order.Total,
			&order.CreateAt,
		)
		order.UserId = userId

		if err != nil {
			return nil, getOrdersByUserIdError
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, rowsIterationError
	}

	return &orders, nil
}

func (r *OrderRepo) GetOrderById(ctx context.Context, orderId int64) (*models2.Order, error) {
	order := new(models2.Order)
	order.Id = orderId
	err := r.db.QueryRow(ctx, getOrderByIdQuery, orderId).Scan(&order.UserId, &order.Status, &order.Total, &order.CreateAt)
	if err != nil {
		return &models2.Order{}, getOrderByIdError
	}

	return order, nil
}
