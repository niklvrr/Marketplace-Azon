package orderService

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type mockRepo struct {
	CreateOrderFn            func(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error)
	GetOrdersByUserIdFn      func(ctx context.Context, userId int64) (*[]model.Order, error)
	GetOrderByIdFn           func(ctx context.Context, orderId int64) (*model.Order, error)
	GetOrderItemsByOrderIdFn func(ctx context.Context, orderId int64) (*[]model.OrderItem, error)
	DeleteOrderByIdFn        func(ctx context.Context, orderId int64) error
}

func (m *mockRepo) CreateOrder(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error) {
	return m.CreateOrderFn(ctx, userId, items)
}
func (m *mockRepo) GetOrdersByUserId(ctx context.Context, userId int64) (*[]model.Order, error) {
	return m.GetOrdersByUserIdFn(ctx, userId)
}
func (m *mockRepo) GetOrderById(ctx context.Context, orderId int64) (*model.Order, error) {
	return m.GetOrderByIdFn(ctx, orderId)
}
func (m *mockRepo) GetOrderItemsByOrderId(ctx context.Context, orderId int64) (*[]model.OrderItem, error) {
	return m.GetOrderItemsByOrderIdFn(ctx, orderId)
}
func (m *mockRepo) DeleteOrderById(ctx context.Context, orderId int64) error {
	return m.DeleteOrderByIdFn(ctx, orderId)
}

func TestOrderService_CreateOrder(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.CreateOrderRequest
		repoFn  func(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error)
		wantID  int64
		wantErr bool
	}{
		{
			"success",
			&model.CreateOrderRequest{
				UserId: 2,
				OrderItems: []model.OrderItemRequest{
					{ProductId: 10, Quantity: 1, Price: 100},
					{ProductId: 11, Quantity: 2, Price: 50},
				},
			},
			func(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error) {
				if userId != 2 {
					t.Fatalf("unexpected userId: got %d want %d", userId, 2)
				}
				if len(*items) != 2 {
					t.Fatalf("unexpected items len: %d", len(*items))
				}
				if (*items)[0].ProductId != 10 || (*items)[0].Quantity != 1 || (*items)[0].Price != 100 {
					t.Fatalf("unexpected first item: %+v", (*items)[0])
				}
				return 77, nil
			},
			77,
			false,
		},
		{
			"repo error",
			&model.CreateOrderRequest{UserId: 3, OrderItems: []model.OrderItemRequest{{ProductId: 1, Quantity: 1, Price: 10}}},
			func(ctx context.Context, userId int64, items *[]model.OrderItem) (int64, error) {
				return 0, errors.New("db")
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{CreateOrderFn: tt.repoFn}
			s := NewOrderService(repo)
			id, err := s.CreateOrder(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if id != tt.wantID {
				t.Fatalf("got id %d want %d", id, tt.wantID)
			}
		})
	}
}

func TestOrderService_GetOrdersByUserId(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.GetOrdersByUserIdRequest
		repoFn  func(ctx context.Context, userId int64) (*[]model.Order, error)
		want    *[]model.OrderResponse
		wantErr bool
	}{
		{
			"success",
			&model.GetOrdersByUserIdRequest{UserId: 5},
			func(ctx context.Context, userId int64) (*[]model.Order, error) {
				ords := []model.Order{
					{Id: 1, UserId: 5, Status: "ok", Total: 100},
					{Id: 2, UserId: 5, Status: "paid", Total: 200},
				}
				return &ords, nil
			},
			&[]model.OrderResponse{{Id: 1, UserId: 5, Status: "ok", Total: 100}, {Id: 2, UserId: 5, Status: "paid", Total: 200}},
			false,
		},
		{
			"repo error",
			&model.GetOrdersByUserIdRequest{UserId: 6},
			func(ctx context.Context, userId int64) (*[]model.Order, error) {
				return nil, errors.New("db")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{GetOrdersByUserIdFn: tt.repoFn}
			s := NewOrderService(repo)
			got, err := s.GetOrdersByUserId(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetOrderById(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.GetOrderByIdRequest
		repoFn  func(ctx context.Context, orderId int64) (*model.Order, error)
		want    *model.OrderResponse
		wantErr bool
	}{
		{
			"success",
			&model.GetOrderByIdRequest{Id: 9},
			func(ctx context.Context, orderId int64) (*model.Order, error) {
				return &model.Order{Id: 9, UserId: 3, Status: "done", Total: 500}, nil
			},
			&model.OrderResponse{Id: 9, UserId: 3, Status: "done", Total: 500},
			false,
		},
		{
			"repo error",
			&model.GetOrderByIdRequest{Id: 10},
			func(ctx context.Context, orderId int64) (*model.Order, error) {
				return nil, errors.New("not found")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{GetOrderByIdFn: tt.repoFn}
			s := NewOrderService(repo)
			got, err := s.GetOrderById(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetOrderItemsByOrderId(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.GetOrderItemsByOrderIdRequest
		repoFn  func(ctx context.Context, orderId int64) (*[]model.OrderItem, error)
		want    *[]model.OrderItemResponse
		wantErr bool
	}{
		{
			"success",
			&model.GetOrderItemsByOrderIdRequest{OrderId: 4},
			func(ctx context.Context, orderId int64) (*[]model.OrderItem, error) {
				items := []model.OrderItem{
					{Id: 1, OrderId: 4, ProductId: 7, Quantity: 2, Price: 50},
					{Id: 2, OrderId: 4, ProductId: 8, Quantity: 1, Price: 100},
				}
				return &items, nil
			},
			&[]model.OrderItemResponse{
				{Id: 1, OrderId: 4, ProductId: 7, Quantity: 2, Price: 50},
				{Id: 2, OrderId: 4, ProductId: 8, Quantity: 1, Price: 100},
			},
			false,
		},
		{
			"repo error",
			&model.GetOrderItemsByOrderIdRequest{OrderId: 5},
			func(ctx context.Context, orderId int64) (*[]model.OrderItem, error) {
				return nil, errors.New("db")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{GetOrderItemsByOrderIdFn: tt.repoFn}
			s := NewOrderService(repo)
			got, err := s.GetOrderItemsByOrderId(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestOrderService_DeleteOrderById(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.DeleteOrderByIdRequest
		repoFn  func(ctx context.Context, orderId int64) error
		wantErr bool
	}{
		{"success", &model.DeleteOrderByIdRequest{Id: 12}, func(ctx context.Context, orderId int64) error { return nil }, false},
		{"repo error", &model.DeleteOrderByIdRequest{Id: 13}, func(ctx context.Context, orderId int64) error { return errors.New("db") }, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{DeleteOrderByIdFn: tt.repoFn}
			s := NewOrderService(repo)
			err := s.DeleteOrderById(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
		})
	}
}
