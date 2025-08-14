package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
	"github.com/niklvrr/myMarketplace/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductsHandler(db *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{svc: service.NewProductService(db)}
}

// TODO func RegisterRoutes
//func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) {
//	rg.GET("/products", h.GetAll)
//	rg.GET("/products/:id", h.Get)
//	rg.POST("/products", RequireRole("seller"), h.Create)
//	rg.PUT("/products/:id", RequireRole("seller"), h.Update)
//	rg.DELETE("/products/:id", RequireRole("seller", "admin"), h.Delete)
//}

func (h *ProductHandler) Create(ctx *gin.Context) {
	var req model.CreateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// TODO RespondError model

		return
	}

	userIdStr, exist := ctx.Get("userId")
	if !exist {
		// TODO RespondIdError model

		return
	}

	userId, ok := strconv.Atoi(userIdStr.(string))
	if ok != nil {
		// TODO EmptyIdError

		return
	}

	product, err := h.svc.Create(ctx, int64(userId), &req)
	if err != nil {
		// TODO RespondServiceError model

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": product})
}

func (h *ProductHandler) Get(ctx *gin.Context) {
	var req model.GetProductsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// TODO RespondError model\

		return
	}

	product, err := h.svc.GetById(ctx, &req)
	if err != nil {
		// TODO RespondServiceError model

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Update(ctx *gin.Context) {
	var req model.UpdateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// TODO RespondError model

		return
	}

	userIdStr, exist := ctx.Get("userId")
	if !exist {
		// TODO RespondIdError model

		return
	}

	userId, ok := strconv.Atoi(userIdStr.(string))
	if ok != nil {
		// TODO EmptyIdError

		return
	}

	product, err := h.svc.UpdateById(ctx, int64(userId), &req)
	if err != nil {
		// TODO RespondServiceError model

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": product})
}

func (h *ProductHandler) Delete(ctx *gin.Context) {
	var req model.DeleteProductRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		// TODO RespondError model

		return
	}

	err := h.svc.DeleteById(ctx, &req)
	if err != nil {
		// TODO RespondServiceError model

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
		// TODO RespondServiceError model

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
		// TODO RespondError model

		return
	}

	products, total, err := h.svc.Search(ctx, page, limit, &req)
	if err != nil {
		// TODO RespondServiceError model

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
