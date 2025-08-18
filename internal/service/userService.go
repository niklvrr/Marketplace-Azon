package service

import (
	"context"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type IUserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserById(ctx context.Context, userId int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUserById(ctx context.Context, user *model.User) error
	BlockUserById(ctx context.Context, userId int64) error
	UnBlockUserById(ctx context.Context, userId int64) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	UpdateUserRole(ctx context.Context, userId int64, newRole string) error
	ApproveProduct(ctx context.Context, productId int64) error
}

type UserService struct {
	repo IUserRepository
}

func NewUserService(repo IUserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SignUp(ctx context.Context, req *model.SighUpRequest) (model.UserResponse, error) {
	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.repo.CreateUser(ctx, &user)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *model.LoginRequest) (model.UserResponse, error) {
	u := model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := s.repo.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error) {
	user, err := s.repo.GetUserById(ctx, req.Id)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) UpdateUserById(ctx context.Context, req *model.UpdateUserByIdRequest) (model.UserResponse, error) {
	user := model.User{
		Id:       req.Id,
		Name:     *req.Name,
		Email:    *req.Email,
		Password: *req.Password,
	}

	err := s.repo.UpdateUserById(ctx, &user)
	if err != nil {
		return model.UserResponse{}, err
	}

	return model.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *UserService) BlockUserById(ctx context.Context, req *model.BlockUserByIdRequest) error {
	err := s.repo.BlockUserById(ctx, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UnblockUserById(ctx context.Context, req *model.UnblockUserByIdRequest) error {
	err := s.repo.UnBlockUserById(ctx, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, req *model.UpdateUserRoleRequest) error {
	err := s.repo.UpdateUserRole(ctx, req.Id, req.Role)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return []model.UserResponse{}, err
	}

	var usersResponse []model.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, model.UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		})
	}

	return usersResponse, nil
}

func (s *UserService) ApproveProduct(ctx context.Context, req *model.ApproveProductRequest) error {
	err := s.repo.ApproveProduct(ctx, req.ProductId)
	if err != nil {
		return err
	}

	return nil
}
