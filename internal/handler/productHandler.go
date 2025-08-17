package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type IProductService interface {
	Create(ctx context.Context, sellerId int64, req *model.CreateProductRequest) (model.ProductResponse, error)
	GetById(ctx context.Context, req *model.GetProductsRequest) (model.ProductResponse, error)
	UpdateById(ctx context.Context, sellerId int64, req *model.UpdateProductRequest) (model.ProductResponse, error)
	DeleteById(ctx context.Context, req *model.DeleteProductRequest) error
	GetAll(ctx context.Context, page, limit int) ([]model.ProductResponse, int64, error)
	Search(ctx context.Context, page, limit int, req *model.SearchProductsRequest) ([]model.ProductResponse, int64, error)
}

type ProductHandler struct {
	svc IProductService
}

func NewProductsHandler(service IProductService) *ProductHandler {
	return &ProductHandler{svc: service}
}

func (h *ProductHandler) Create(ctx *gin.Context) {
	var req model.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	userIdStr, exist := ctx.Get("userId")
	if !exist {
		errs.RespondError(ctx, http.StatusUnauthorized, "unauthorized", "user is not authorized")
		return
	}

	userId, ok := strconv.Atoi(userIdStr.(string))
	if ok != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "validation_error", "invalid user id")
		return
	}

	product, err := h.svc.Create(ctx, int64(userId), &req)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *ProductHandler) Get(ctx *gin.Context) {
	var req model.GetProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	product, err := h.svc.GetById(ctx, &req)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Update(ctx *gin.Context) {
	var req model.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	userIdStr, exist := ctx.Get("userId")
	if !exist {
		errs.RespondError(ctx, http.StatusUnauthorized, "unauthorized", "user is not authorized")
		return
	}

	userId, ok := strconv.Atoi(userIdStr.(string))
	if ok != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "validation_error", "invalid user id")
		return
	}

	product, err := h.svc.UpdateById(ctx, int64(userId), &req)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Delete(ctx *gin.Context) {
	var req model.DeleteProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.DeleteById(ctx, &req)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "product deleted"})
}

func (h *ProductHandler) GetAll(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	products, total, err := h.svc.GetAll(ctx, page, limit)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":       products,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}

func (h *ProductHandler) Search(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	var req model.SearchProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		errs.RespondError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	products, total, err := h.svc.Search(ctx, page, limit, &req)
	if err != nil {
		errs.RespondServiceError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data":       products,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": (total + int64(limit) - 1) / int64(limit),
	})
}
