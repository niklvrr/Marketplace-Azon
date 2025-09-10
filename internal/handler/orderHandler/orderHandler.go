package orderHandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (int64, error)
	GetOrdersByUserId(ctx context.Context, req *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error)
	GetOrderById(ctx context.Context, req *model.GetOrderByIdRequest) (*model.OrderResponse, error)
	GetOrderItemsByOrderId(ctx context.Context, req *model.GetOrderItemsByOrderIdRequest) (*[]model.OrderItemResponse, error)
	DeleteOrderById(ctx context.Context, req *model.DeleteOrderByIdRequest) error
}

type OrderHandler struct {
	svc IOrderService
}

func NewOrderHandler(svc IOrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	orderId, err := h.svc.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": orderId})
}

func (h *OrderHandler) GetOrdersByUserId(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		errs.RespondError(c, http.StatusUnauthorized, "unauthorized", "user id not found")
		return
	}

	req := &model.GetOrdersByUserIdRequest{UserId: userId.(int64)}
	resp, err := h.svc.GetOrdersByUserId(c.Request.Context(), req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": resp})
}

func (h *OrderHandler) GetOrderById(c *gin.Context) {
	orderId := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.GetOrderByIdRequest{Id: int64(orderIdInt)}
	order, err := h.svc.GetOrderById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (h *OrderHandler) GetOrderItemsByOrderId(c *gin.Context) {
	orderId := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.GetOrderItemsByOrderIdRequest{OrderId: int64(orderIdInt)}
	resp, err := h.svc.GetOrderItemsByOrderId(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"order_items": resp})
}

func (h *OrderHandler) DeleteOrderById(c *gin.Context) {
	orderId := c.Param("id")
	orderIdInt, err := strconv.Atoi(orderId)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.DeleteOrderByIdRequest{Id: int64(orderIdInt)}
	err = h.svc.DeleteOrderById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}
