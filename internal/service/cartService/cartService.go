package cartService

import (
	"context"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type ICartRepository interface {
	GetCartByUserId(ctx context.Context, userId int64) (*model.Cart, error)
	GetCartItemsByCartId(ctx context.Context, cartId int64) (*[]model.CartItem, error)
	AddItem(ctx context.Context, cartId, productId int64, quantity int) (int64, error)
	RemoveItem(ctx context.Context, cartId, id int64) error
	ClearCart(ctx context.Context, cartId int64) error
}

type CartService struct {
	repo ICartRepository
}

func NewCartService(repo ICartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCartByUserId(ctx context.Context, req *model.GetCartByUserIdRequest) (*model.CartResponse, error) {
	cart, err := s.repo.GetCartByUserId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &model.CartResponse{
		Id: cart.Id,
	}, nil
}

func (s *CartService) GetCartItemsByCartId(ctx context.Context, req *model.GetCartItemsByCartIdRequest) (*[]model.CartItemResponse, error) {
	cartId := req.CartId
	cartItems, err := s.repo.GetCartItemsByCartId(ctx, cartId)
	if err != nil {
		return nil, err
	}

	var items []model.CartItemResponse
	for _, item := range *cartItems {
		items = append(items, model.CartItemResponse{
			Id:        item.Id,
			CartId:    cartId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	return &items, nil
}
func (s *CartService) AddItem(ctx context.Context, req *model.AddItemRequest) (int64, error) {
	itemId, err := s.repo.AddItem(ctx, req.CartId, req.ProductId, req.Quantity)
	if err != nil {
		return 0, err
	}

	return itemId, nil
}

func (s *CartService) RemoveItem(ctx context.Context, req *model.RemoveItemRequest) error {
	err := s.repo.RemoveItem(ctx, req.CartId, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CartService) ClearCart(ctx context.Context, req *model.ClearCartRequest) error {
	err := s.repo.ClearCart(ctx, req.CartId)
	if err != nil {
		return err
	}

	return nil
}
