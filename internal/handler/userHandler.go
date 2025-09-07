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
	SignUp(ctx context.Context, req *model.SighUpRequest) (string, error)
	Login(ctx context.Context, req *model.LoginRequest) (string, error)
	GetUserById(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error)
	UpdateUserById(ctx context.Context, req *model.UpdateUserByIdRequest, role string) (model.UserResponse, error)
	GetUserByEmail(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error)
	BlockUserById(ctx context.Context, req *model.BlockUserByIdRequest) error
	UnblockUserById(ctx context.Context, req *model.UnblockUserByIdRequest) error
	GetAllUsers(ctx context.Context) ([]model.UserResponse, error)
	UpdateUserRole(ctx context.Context, req *model.UpdateUserRoleRequest) error
	ApproveProduct(ctx context.Context, req *model.ApproveProductRequest) error
	Logout(ctx context.Context, req *model.LogoutRequest) error
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

	token, err := h.svc.SignUp(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": token})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	token, err := h.svc.Login(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id, exist := c.Get("user_id")
	if !exist {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no user found")
		return
	}
	req := model.GetUserByIdRequest{Id: id.(int64)}

	user, err := h.svc.GetUserById(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	var req model.GetUserByEmailRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	user, err := h.svc.GetUserByEmail(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) UpdateUserById(c *gin.Context) {
	id, exist := c.Get("user_id")
	if !exist {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no user found")
		return
	}

	role, exist := c.Get("role")
	if !exist {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no role found")
		return
	}

	var req model.UpdateUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	req.Id = id.(int64)

	user, err := h.svc.UpdateUserById(c, &req, role.(string))
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *UserHandler) Logout(c *gin.Context) {
	id, exist := c.Get("user_id")
	if !exist {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", "no user found")
		return
	}

	blacklistKey := "blacklist_user:" + strconv.Itoa(int(id.(int64)))
	req := model.LogoutRequest{BlockKey: blacklistKey}
	err := h.svc.Logout(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": true})
}

func (h *UserHandler) BlockUserById(c *gin.Context) {
	var blockReq model.BlockUserByIdRequest
	if err := c.ShouldBind(&blockReq); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.BlockUserById(c, &blockReq)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	blockKey := "blocked_user:" + strconv.Itoa(int(blockReq.Id))
	logoutReq := model.LogoutRequest{BlockKey: blockKey}
	err = h.svc.Logout(c, &logoutReq)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *UserHandler) UnblockUserById(c *gin.Context) {
	var req model.UnblockUserByIdRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.UnblockUserById(c, &req)
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
	var req model.UpdateUserRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.UpdateUserRole(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}

func (h *UserHandler) ApproveProduct(c *gin.Context) {
	var req model.ApproveProductRequest
	if err := c.ShouldBind(&req); err != nil {
		errs.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err := h.svc.ApproveProduct(c, &req)
	if err != nil {
		errs.RespondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
