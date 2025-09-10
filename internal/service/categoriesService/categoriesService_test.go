package categoriesService

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/niklvrr/myMarketplace/internal/model"
)

type mockRepo struct {
	CreateCategoryFn   func(ctx context.Context, category *model.Category) error
	GetCategoryByIdFn  func(ctx context.Context, id int64) (*model.Category, error)
	UpdateCategoryFn   func(ctx context.Context, category *model.Category) error
	DeleteCategoryFn   func(ctx context.Context, id int64) error
	GetAllCategoriesFn func(ctx context.Context) (*[]model.Category, error)
}

func (m *mockRepo) CreateCategory(ctx context.Context, category *model.Category) error {
	return m.CreateCategoryFn(ctx, category)
}
func (m *mockRepo) GetCategoryById(ctx context.Context, id int64) (*model.Category, error) {
	return m.GetCategoryByIdFn(ctx, id)
}
func (m *mockRepo) UpdateCategory(ctx context.Context, category *model.Category) error {
	return m.UpdateCategoryFn(ctx, category)
}
func (m *mockRepo) DeleteCategory(ctx context.Context, id int64) error {
	return m.DeleteCategoryFn(ctx, id)
}
func (m *mockRepo) GetAllCategories(ctx context.Context) (*[]model.Category, error) {
	return m.GetAllCategoriesFn(ctx)
}

func TestCategoriesService_Create(t *testing.T) {
	tests := []struct {
		name     string
		req      *model.CreateCategoryRequest
		repoFn   func(ctx context.Context, category *model.Category) error
		wantResp *model.CategoryResponse
		wantErr  bool
	}{
		{
			"success",
			&model.CreateCategoryRequest{Name: "Books", Description: "All books"},
			func(ctx context.Context, category *model.Category) error {
				category.Id = 12
				return nil
			},
			&model.CategoryResponse{Id: 12, Name: "Books", Description: "All books"},
			false,
		},
		{
			"repo error",
			&model.CreateCategoryRequest{Name: "X", Description: "Y"},
			func(ctx context.Context, category *model.Category) error {
				return errors.New("repo error")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{
				CreateCategoryFn: tt.repoFn,
			}
			s := NewCategoriesService(repo)
			got, err := s.Create(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Fatalf("got %+v want %+v", got, tt.wantResp)
			}
		})
	}
}

func TestCategoriesService_GetById(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.GetCategoryByIdRequest
		repoFn  func(ctx context.Context, id int64) (*model.Category, error)
		want    *model.CategoryResponse
		wantErr bool
	}{
		{
			"success",
			&model.GetCategoryByIdRequest{Id: 5},
			func(ctx context.Context, id int64) (*model.Category, error) {
				return &model.Category{Id: id, Name: "C", Description: "D"}, nil
			},
			&model.CategoryResponse{Id: 5, Name: "C", Description: "D"},
			false,
		},
		{
			"repo error",
			&model.GetCategoryByIdRequest{Id: 6},
			func(ctx context.Context, id int64) (*model.Category, error) {
				return nil, errors.New("not found")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{GetCategoryByIdFn: tt.repoFn}
			s := NewCategoriesService(repo)
			got, err := s.GetById(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestCategoriesService_Update(t *testing.T) {
	name := "Updated"
	desc := "NewDesc"
	tests := []struct {
		name    string
		req     *model.UpdateCategoryRequest
		repoFn  func(ctx context.Context, category *model.Category) error
		want    *model.CategoryResponse
		wantErr bool
	}{
		{
			"success",
			&model.UpdateCategoryRequest{Id: 7, Name: &name, Description: &desc},
			func(ctx context.Context, category *model.Category) error {
				if category.Id != 7 || category.Name != name || category.Description != desc {
					return errors.New("invalid category passed")
				}
				return nil
			},
			&model.CategoryResponse{Id: 7, Name: name, Description: desc},
			false,
		},
		{
			"repo error",
			&model.UpdateCategoryRequest{Id: 8, Name: &name, Description: &desc},
			func(ctx context.Context, category *model.Category) error {
				return errors.New("update failed")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{UpdateCategoryFn: tt.repoFn}
			s := NewCategoriesService(repo)
			got, err := s.Update(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}

func TestCategoriesService_Delete(t *testing.T) {
	tests := []struct {
		name    string
		req     *model.DeleteCategoryRequest
		repoFn  func(ctx context.Context, id int64) error
		wantErr bool
	}{
		{"success", &model.DeleteCategoryRequest{Id: 9}, func(ctx context.Context, id int64) error { return nil }, false},
		{"repo error", &model.DeleteCategoryRequest{Id: 10}, func(ctx context.Context, id int64) error { return errors.New("del") }, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{DeleteCategoryFn: tt.repoFn}
			s := NewCategoriesService(repo)
			err := s.Delete(context.Background(), tt.req)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
		})
	}
}

func TestCategoriesService_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		repoFn  func(ctx context.Context) (*[]model.Category, error)
		want    *[]model.CategoryResponse
		wantErr bool
	}{
		{
			"success",
			func(ctx context.Context) (*[]model.Category, error) {
				cs := []model.Category{{Id: 1, Name: "A", Description: "a"}, {Id: 2, Name: "B", Description: "b"}}
				return &cs, nil
			},
			&[]model.CategoryResponse{{Id: 1, Name: "A", Description: "a"}, {Id: 2, Name: "B", Description: "b"}},
			false,
		},
		{
			"repo error",
			func(ctx context.Context) (*[]model.Category, error) {
				return nil, errors.New("repo err")
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepo{GetAllCategoriesFn: tt.repoFn}
			s := NewCategoriesService(repo)
			got, err := s.GetAll(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v want %+v", got, tt.want)
			}
		})
	}
}
