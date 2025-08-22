package service

import (
	"context"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type IOrderRepository interface {
	CreateOrder(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error)
	GetOrdersByUserId(ctx context.Context, userId int64) (*[]model.Order, error)
	GetOrderById(ctx context.Context, orderId int64) (*model.Order, error)
	GetOrderItemsByOrderId(ctx context.Context, orderId int64) (*[]model.OrderItem, error)
	DeleteOrderById(ctx context.Context, orderId int64) error
}

type OrderService struct {
	repo IOrderRepository
}

func NewOrderService(repo IOrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (int64, error) {
	var items []model.OrderItem
	for _, r := range req.OrderItems {
		items = append(items, model.OrderItem{
			ProductId: r.ProductId,
			Quantity:  r.Quantity,
			Price:     r.Price,
		})
	}

	orderId, err := s.repo.CreateOrder(ctx, req.UserId, &items)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (s *OrderService) GetOrdersByUserId(ctx context.Context, req *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
	orders, err := s.repo.GetOrdersByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	var resp []model.OrderResponse
	for _, o := range *orders {
		resp = append(resp, model.OrderResponse{
			Id:     o.Id,
			UserId: o.UserId,
			Status: o.Status,
			Total:  o.Total,
		})
	}
	return &resp, nil
}

func (s *OrderService) GetOrderById(ctx context.Context, req *model.GetOrderByIdRequest) (*model.OrderResponse, error) {
	order, err := s.repo.GetOrderById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &model.OrderResponse{
		Id:     order.Id,
		UserId: order.UserId,
		Status: order.Status,
		Total:  order.Total,
	}, nil
}

func (s *OrderService) GetOrderItemsByOrderId(ctx context.Context, req *model.GetOrderItemsByOrderIdRequest) (*[]model.OrderItemResponse, error) {
	items, err := s.repo.GetOrderItemsByOrderId(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	var resp []model.OrderItemResponse
	for _, o := range *items {
		resp = append(resp, model.OrderItemResponse{
			Id:        o.Id,
			OrderId:   o.OrderId,
			ProductId: o.ProductId,
			Quantity:  o.Quantity,
			Price:     o.Price,
		})
	}

	return &resp, nil
}

func (s *OrderService) DeleteOrderById(ctx context.Context, req *model.DeleteOrderByIdRequest) error {
	err := s.repo.DeleteOrderById(ctx, req.Id)
	if err != nil {
		return err
	}

	return nil
}
