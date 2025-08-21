package service

import (
	"context"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type ICategoriesRepository interface {
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategoryById(ctx context.Context, id int64) (*model.Category, error)
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	GetAllCategories(ctx context.Context) (*[]model.Category, error)
}

type CategoriesService struct {
	repo ICategoriesRepository
}

func NewCategoriesService(repo ICategoriesRepository) *CategoriesService {
	return &CategoriesService{repo: repo}
}

func (s *CategoriesService) Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	c := model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.CreateCategory(ctx, &c); err != nil {
		return nil, err
	}

	return &model.CategoryResponse{
		Id:          c.Id,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (s *CategoriesService) GetById(ctx context.Context, req *model.GetCategoryByIdRequest) (*model.CategoryResponse, error) {
	c, err := s.repo.GetCategoryById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &model.CategoryResponse{
		Id:          c.Id,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (s *CategoriesService) Update(ctx context.Context, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	c := model.Category{
		Id:          req.Id,
		Name:        *req.Name,
		Description: *req.Description,
	}

	if err := s.repo.UpdateCategory(ctx, &c); err != nil {
		return nil, err
	}

	return &model.CategoryResponse{
		Id:          c.Id,
		Name:        c.Name,
		Description: c.Description,
	}, nil
}

func (s *CategoriesService) Delete(ctx context.Context, req *model.DeleteCategoryRequest) error {
	if err := s.repo.DeleteCategory(ctx, req.Id); err != nil {
		return err
	}

	return nil
}

func (s *CategoriesService) GetAll(ctx context.Context) (*[]model.CategoryResponse, error) {
	cs, err := s.repo.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	var resp []model.CategoryResponse
	for _, c := range *cs {
		resp = append(resp, model.CategoryResponse{
			Id:          c.Id,
			Name:        c.Name,
			Description: c.Description,
		})
	}

	return &resp, nil
}
