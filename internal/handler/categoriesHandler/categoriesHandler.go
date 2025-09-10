package categoriesHandler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type ICategoriesService interface {
	Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error)
	GetById(ctx context.Context, req *model.GetCategoryByIdRequest) (*model.CategoryResponse, error)
	Update(ctx context.Context, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error)
	Delete(ctx context.Context, req *model.DeleteCategoryRequest) error
	GetAll(ctx context.Context) (*[]model.CategoryResponse, error)
}

type CategoriesHandler struct {
	svc ICategoriesService
}

func NewCategoryHandler(svc ICategoriesService) *CategoriesHandler {
	return &CategoriesHandler{svc: svc}
}

func (h *CategoriesHandler) Create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	cat, err := h.svc.Create(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": cat})
}

func (h *CategoriesHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.GetCategoryByIdRequest{
		Id: int64(idInt),
	}

	cat, err := h.svc.GetById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func (h *CategoriesHandler) Update(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.Id = int64(idInt)

	cat, err := h.svc.Update(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cat})
}

func (h *CategoriesHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	req := model.DeleteCategoryRequest{Id: int64(idInt)}

	err = h.svc.Delete(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *CategoriesHandler) GetAll(c *gin.Context) {
	cats, err := h.svc.GetAll(c)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cats})
}
