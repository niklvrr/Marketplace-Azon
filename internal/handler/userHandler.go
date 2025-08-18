package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/errs"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type IUserService interface {
	SignUp(ctx context.Context, req *model.SighUpRequest) (model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginRequest) (model.UserResponse, error)
	GetUserById(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error)
	UpdateUserById(ctx context.Context, req *model.UpdateUserByIdRequest) (model.UserResponse, error)
	GetUserByEmail(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error)
	BlockUserById(ctx context.Context, req *model.BlockUserByIdRequest) error
	UnblockUserById(ctx context.Context, req *model.UnblockUserByIdRequest) error
	GetAllUsers(ctx context.Context) ([]model.UserResponse, error)
	UpdateUserRole(ctx context.Context, req *model.UpdateUserRoleRequest) error
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
	email := c.Query("email")
	if email == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "email is required")
		return
	}

	password := c.Query("password")
	if password == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "password is required")
		return
	}

	req := model.LoginRequest{Email: email, Password: password}

	user, err := h.svc.Login(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	req := model.GetUserByIdRequest{Id: int64(idInt)}

	user, err := h.svc.GetUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no email provided")
		return
	}
	req := model.GetUserByEmailRequest{Email: email}

	user, err := h.svc.GetUserByEmail(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no id provided")
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var req model.UpdateUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.Id = int64(idInt)

	user, err := h.svc.UpdateUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) BlockUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no id provided")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req := model.BlockUserByIdRequest{Id: int64(idInt)}

	err = h.svc.BlockUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *UserHandler) UnblockUserById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no id provided")
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req := model.UnblockUserByIdRequest{Id: int64(idInt)}

	err = h.svc.UnblockUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers(c)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no id provided")
		return
	}
	idInt, err := strconv.Atoi(id)
	req := model.UpdateUserRoleRequest{Id: int64(idInt)}

	role := c.Query("role")
	if role == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no role provided")
		return
	}
	req.Role = role

	err = h.svc.UpdateUserRole(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *UserHandler) ApproveProduct(c *gin.Context) {
	productId := c.Param("productId")
	if productId == "" {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no productId provided")
		return
	}
	idInt, err := strconv.Atoi(productId)
	req := model.ApproveProductRequest{ProductId: int64(idInt)}

	err = h.svc.ApproveProduct(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
