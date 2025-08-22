package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type ICartService interface {
	GetCartByUserId(ctx context.Context, req *model.GetCartByUserIdRequest) (*model.CartResponse, error)
	GetCartItemsByCartId(ctx context.Context, req *model.GetCartItemsByCartIdRequest) (*[]model.CartItemResponse, error)
	AddItem(ctx context.Context, req *model.AddItemRequest) (int64, error)
	RemoveItem(ctx context.Context, req *model.RemoveItemRequest) error
	ClearCart(ctx context.Context, req *model.ClearCartRequest) error
}

type CartHandler struct {
	svc ICartService
}

func NewCartHandler(svc ICartService) *CartHandler {
	return &CartHandler{svc: svc}
}

func (h *CartHandler) GetCartByUserId(c *gin.Context) {
	userId, exist := c.Get("user_id")
	if !exist {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "user id not found")
		return
	}

	req := model.GetCartByUserIdRequest{UserId: userId.(int64)}
	cart, err := h.svc.GetCartByUserId(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) GetCartItemsByCartId(c *gin.Context) {
	cartId := c.Param("cart_id")
	cartIdInt, err := strconv.Atoi(cartId)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.GetCartItemsByCartIdRequest{CartId: int64(cartIdInt)}
	cart, err := h.svc.GetCartItemsByCartId(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

func (h *CartHandler) AddItem(c *gin.Context) {
	var req model.AddItemRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	itemId, err := h.svc.AddItem(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart_item_id": itemId})
}

func (h *CartHandler) RemoveItem(c *gin.Context) {
	var req model.RemoveItemRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.RemoveItem(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	var req model.ClearCartRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.ClearCart(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
