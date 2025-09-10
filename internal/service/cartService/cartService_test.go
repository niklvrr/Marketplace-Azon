package cartService

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type mockRepo struct {
	GetCartByUserIdFn      func(ctx context.Context, userId int64) (*model.Cart, error)
	GetCartItemsByCartIdFn func(ctx context.Context, cartId int64) (*[]model.CartItem, error)
	AddItemFn              func(ctx context.Context, cartId, productId int64, quantity int) (int64, error)
	RemoveItemFn           func(ctx context.Context, cartId, id int64) error
	ClearCartFn            func(ctx context.Context, cartId int64) error
}

func (m *mockRepo) GetCartByUserId(ctx context.Context, userId int64) (*model.Cart, error) {
	return m.GetCartByUserIdFn(ctx, userId)
}
func (m *mockRepo) GetCartItemsByCartId(ctx context.Context, cartId int64) (*[]model.CartItem, error) {
	return m.GetCartItemsByCartIdFn(ctx, cartId)
}
func (m *mockRepo) AddItem(ctx context.Context, cartId, productId int64, quantity int) (int64, error) {
	return m.AddItemFn(ctx, cartId, productId, quantity)
}
func (m *mockRepo) RemoveItem(ctx context.Context, cartId, id int64) error {
	return m.RemoveItemFn(ctx, cartId, id)
}
func (m *mockRepo) ClearCart(ctx context.Context, cartId int64) error {
	return m.ClearCartFn(ctx, cartId)
}

func TestCartService_GetCartByUserId(t *testing.T) {
	tests := []struct {
		name       string
		userId     int64
		repoCart   *model.Cart
		repoErr    error
		wantErr    bool
		wantCartID int64
	}{
		{"success", 10, &model.Cart{Id: 100}, nil, false, 100},
		{"repo error", 11, nil, errors.New("repo error"), true, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				GetCartByUserIdFn: func(ctx context.Context, userId int64) (*model.Cart, error) {
					if userId != tt.userId {
						t.Fatalf("unexpected userId: got %d want %d", userId, tt.userId)
					}
					return tt.repoCart, tt.repoErr
				},
			}
			s := NewCartService(repo)
			resp, err := s.GetCartByUserId(context.Background(), &model.GetCartByUserIdRequest{UserId: tt.userId})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if resp.Id != tt.wantCartID {
				t.Fatalf("got id %d want %d", resp.Id, tt.wantCartID)
			}
		})
	}
}

func TestCartService_GetCartItemsByCartId(t *testing.T) {
	tests := []struct {
		name      string
		cartId    int64
		repoItems *[]model.CartItem
		repoErr   error
		wantErr   bool
		wantItems []model.CartItemResponse
	}{
		{
			"success",
			5,
			&[]model.CartItem{
				{Id: 1, ProductId: 11, Quantity: 2},
				{Id: 2, ProductId: 22, Quantity: 3},
			},
			nil,
			false,
			[]model.CartItemResponse{
				{Id: 1, CartId: 5, ProductId: 11, Quantity: 2},
				{Id: 2, CartId: 5, ProductId: 22, Quantity: 3},
			},
		},
		{"repo error", 6, nil, errors.New("repo err"), true, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				GetCartItemsByCartIdFn: func(ctx context.Context, cartId int64) (*[]model.CartItem, error) {
					if cartId != tt.cartId {
						t.Fatalf("unexpected cartId: got %d want %d", cartId, tt.cartId)
					}
					return tt.repoItems, tt.repoErr
				},
			}
			s := NewCartService(repo)
			resp, err := s.GetCartItemsByCartId(context.Background(), &model.GetCartItemsByCartIdRequest{CartId: tt.cartId})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(*resp, tt.wantItems) {
				t.Fatalf("got items %+v want %+v", *resp, tt.wantItems)
			}
		})
	}
}

func TestCartService_AddItem(t *testing.T) {
	tests := []struct {
		name       string
		req        *model.AddItemRequest
		repoID     int64
		repoErr    error
		wantErr    bool
		wantItemID int64
	}{
		{"success", &model.AddItemRequest{CartId: 7, ProductId: 8, Quantity: 2}, 55, nil, false, 55},
		{"repo error", &model.AddItemRequest{CartId: 7, ProductId: 8, Quantity: 2}, 0, errors.New("repo err"), true, 0},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				AddItemFn: func(ctx context.Context, cartId, productId int64, quantity int) (int64, error) {
					if cartId != tt.req.CartId || productId != tt.req.ProductId || quantity != tt.req.Quantity {
						t.Fatalf("unexpected args")
					}
					return tt.repoID, tt.repoErr
				},
			}
			s := NewCartService(repo)
			id, err := s.AddItem(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if id != tt.wantItemID {
				t.Fatalf("got id %d want %d", id, tt.wantItemID)
			}
		})
	}
}

func TestCartService_RemoveItem(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.RemoveItemRequest
		repoErr error
		wantErr bool
	}{
		{"success", &model.RemoveItemRequest{CartId: 3, Id: 4}, nil, false},
		{"repo error", &model.RemoveItemRequest{CartId: 3, Id: 4}, errors.New("repo err"), true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				RemoveItemFn: func(ctx context.Context, cartId, id int64) error {
					if cartId != tt.req.CartId || id != tt.req.Id {
						t.Fatalf("unexpected args")
					}
					return tt.repoErr
				},
			}
			s := NewCartService(repo)
			err := s.RemoveItem(context.Background(), tt.req)
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

func TestCartService_ClearCart(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.ClearCartRequest
		repoErr error
		wantErr bool
	}{
		{"success", &model.ClearCartRequest{CartId: 9}, nil, false},
		{"repo error", &model.ClearCartRequest{CartId: 9}, errors.New("repo err"), true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				ClearCartFn: func(ctx context.Context, cartId int64) error {
					if cartId != tt.req.CartId {
						t.Fatalf("unexpected cartId")
					}
					return tt.repoErr
				},
			}
			s := NewCartService(repo)
			err := s.ClearCart(context.Background(), tt.req)
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
