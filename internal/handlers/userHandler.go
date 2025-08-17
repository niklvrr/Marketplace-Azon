package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type IUserService interface {
	SignUp(ctx context.Context, req *model.SighUpRequest) (model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (model.UserResponse, error)
	GetUserById(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error)
	UpdateUserById(ctx context.Context, req *model.UpdateUserByIdRequest) (model.UserResponse, error)
	BlockUserById(ctx context.Context, req *model.BlockUserByIdRequest) (model.UserResponse, error)
	UnblockUserById(ctx context.Context, req *model.UnblockUserByIdRequest) (model.UserResponse, error)
	GetAllUsers(ctx context.Context) ([]model.UserResponse, error)
	UpdateUserRole(ctx context.Context, req *model.UpdateUserRoleRequest) (model.UserResponse, error)
	ApproveProduct(ctx context.Context, req *model.ApproveProductRequest) error
}

type UserHandler struct {
	svc IUserService
}

func NewUserHandler(svc IUserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var req model.SighUpRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.svc.SignUp(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.svc.Login(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	var req model.GetUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.svc.GetUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	var req model.UpdateUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.svc.UpdateUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) BlockUserById(c *gin.Context) {
	var req model.BlockUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
}
