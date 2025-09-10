package productService

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v9"
	"github.com/niklvrr/myMarketplace/internal/model"
	"github.com/redis/go-redis/v9"
)

type mockRepo struct {
	CreateProductFn     func(ctx context.Context, product *model.Product) error
	GetProductByIdFn    func(ctx context.Context, productId int64) (*model.Product, error)
	UpdateProductByIdFn func(ctx context.Context, product *model.Product) error
	DeleteProductByIdFn func(ctx context.Context, productId int64) error
	GetAllProductsFn    func(ctx context.Context, offset, limit int) (*[]model.Product, int64, error)
	SearchProductsFn    func(ctx context.Context, text *string, categoryId *int64, min, max *float64, offset, limit int) (*[]model.Product, int64, error)
}

func (m *mockRepo) CreateProduct(ctx context.Context, product *model.Product) error {
	return m.CreateProductFn(ctx, product)
}
func (m *mockRepo) GetProductById(ctx context.Context, productId int64) (*model.Product, error) {
	return m.GetProductByIdFn(ctx, productId)
}
func (m *mockRepo) UpdateProductById(ctx context.Context, product *model.Product) error {
	return m.UpdateProductByIdFn(ctx, product)
}
func (m *mockRepo) DeleteProductById(ctx context.Context, productId int64) error {
	return m.DeleteProductByIdFn(ctx, productId)
}
func (m *mockRepo) GetAllProducts(ctx context.Context, offset, limit int) (*[]model.Product, int64, error) {
	return m.GetAllProductsFn(ctx, offset, limit)
}
func (m *mockRepo) SearchProducts(ctx context.Context, text *string, categoryId *int64, min, max *float64, offset, limit int) (*[]model.Product, int64, error) {
	return m.SearchProductsFn(ctx, text, categoryId, min, max, offset, limit)
}

func TestProductService_Create(t *testing.T) {
	repo := &mockRepo{
		CreateProductFn: func(ctx context.Context, product *model.Product) error {
			product.Id = 21
			return nil
		},
	}
	client, mock := redismock.NewClientMock()
	mock.ExpectDel("products:all").SetVal(1)
	s := NewProductService(repo, client)
	req := &model.CreateProductRequest{
		CategoryId:  2,
		Name:        "P",
		Description: "D",
		Price:       100,
		Stock:       5,
	}
	out, err := s.Create(context.Background(), 9, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if out.Id != 21 || out.SellerId != 9 || out.CategoryId != 2 || out.Name != "P" {
		t.Fatalf("unexpected out: %+v", out)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestProductService_GetById(t *testing.T) {
	repo := &mockRepo{
		GetProductByIdFn: func(ctx context.Context, productId int64) (*model.Product, error) {
			if productId == 5 {
				return &model.Product{Id: 5, SellerId: 2, CategoryId: 3, Name: "X", Price: 10, Stock: 1}, nil
			}
			return nil, errors.New("not found")
		},
	}
	client, _ := redismock.NewClientMock()
	s := NewProductService(repo, client)
	got, err := s.GetById(context.Background(), &model.GetProductsRequest{Id: 5})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.Id != 5 || got.Name != "X" {
		t.Fatalf("unexpected got: %+v", got)
	}
	_, err = s.GetById(context.Background(), &model.GetProductsRequest{Id: 7})
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProductService_UpdateById(t *testing.T) {
	repo := &mockRepo{
		UpdateProductByIdFn: func(ctx context.Context, product *model.Product) error {
			product.Id = 33
			return nil
		},
	}
	client, mock := redismock.NewClientMock()
	mock.ExpectDel("products:all").SetVal(1)
	s := NewProductService(repo, client)
	cat := int64(3)
	name := "N"
	desc := "D"
	price := 150.0
	stock := 7
	req := &model.UpdateProductRequest{
		CategoryId:  &cat,
		Name:        &name,
		Description: &desc,
		Price:       &price,
		Stock:       &stock,
	}
	got, err := s.UpdateById(context.Background(), 12, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.Id != 33 {
		t.Fatalf("unexpected id: %d", got.Id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestProductService_DeleteById(t *testing.T) {
	repo := &mockRepo{
		DeleteProductByIdFn: func(ctx context.Context, productId int64) error {
			if productId == 4 {
				return nil
			}
			return errors.New("db")
		},
	}
	client, mock := redismock.NewClientMock()
	mock.ExpectDel("products:all").SetVal(1)
	s := NewProductService(repo, client)
	if err := s.DeleteById(context.Background(), &model.DeleteProductRequest{Id: 4}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := s.DeleteById(context.Background(), &model.DeleteProductRequest{Id: 5}); err == nil {
		t.Fatalf("expected error")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}
}

func TestProductService_GetAll_CacheHitAndMiss(t *testing.T) {
	products := []model.ProductResponse{
		{Id: 1, SellerId: 1, CategoryId: 2, Name: "A", Price: 10, Stock: 5},
		{Id: 2, SellerId: 1, CategoryId: 2, Name: "B", Price: 20, Stock: 3},
	}
	clientHit, mockHit := redismock.NewClientMock()
	data, _ := json.Marshal(products)
	mockHit.ExpectGet("products:all").SetVal(string(data))
	sHit := NewProductService(nil, clientHit)
	got, total, err := sHit.GetAll(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if total != int64(len(products)) || !reflect.DeepEqual(got, products) {
		t.Fatalf("unexpected cached result got %+v %d", got, total)
	}
	if err := mockHit.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}

	repo := &mockRepo{
		GetAllProductsFn: func(ctx context.Context, offset, limit int) (*[]model.Product, int64, error) {
			prod := []model.Product{
				{Id: 3, SellerId: 2, CategoryId: 4, Name: "C", Price: 30, Stock: 2},
			}
			return &prod, 1, nil
		},
	}
	clientMiss, mockMiss := redismock.NewClientMock()
	mockMiss.ExpectGet("products:all").SetErr(redis.Nil)
	expectedResult := []model.ProductResponse{{Id: 3, SellerId: 2, CategoryId: 4, Name: "C", Price: 30, Stock: 2}}
	dataToCache, _ := json.Marshal(expectedResult)
	mockMiss.ExpectSet("products:all", string(dataToCache), 5*time.Minute).SetVal("OK")
	sMiss := NewProductService(repo, clientMiss)
	got2, total2, err := sMiss.GetAll(context.Background(), 1, 20)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if total2 != 1 || !reflect.DeepEqual(got2, expectedResult) {
		t.Fatalf("unexpected result got %+v %d", got2, total2)
	}
	if err := mockMiss.ExpectationsWereMet(); err != nil {
		t.Fatalf("redis expectations: %v", err)
	}

	repoErr := &mockRepo{
		GetAllProductsFn: func(ctx context.Context, offset, limit int) (*[]model.Product, int64, error) {
			return nil, 0, errors.New("db")
		},
	}
	clientErr, _ := redismock.NewClientMock()
	sErr := NewProductService(repoErr, clientErr)
	_, _, err = sErr.GetAll(context.Background(), 1, 20)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProductService_Search(t *testing.T) {
	repo := &mockRepo{
		SearchProductsFn: func(ctx context.Context, text *string, categoryId *int64, min, max *float64, offset, limit int) (*[]model.Product, int64, error) {
			prod := []model.Product{
				{Id: 7, SellerId: 3, CategoryId: 5, Name: "S", Price: 99, Stock: 1},
			}
			return &prod, 1, nil
		},
	}
	client, _ := redismock.NewClientMock()
	s := NewProductService(repo, client)
	text := "q"
	req := &model.SearchProductsRequest{Text: &text}
	got, total, err := s.Search(context.Background(), 1, 10, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if total != 1 || len(got) != 1 || got[0].Id != 7 {
		t.Fatalf("unexpected result: %+v %d", got, total)
	}

	repoErr := &mockRepo{
		SearchProductsFn: func(ctx context.Context, text *string, categoryId *int64, min, max *float64, offset, limit int) (*[]model.Product, int64, error) {
			return nil, 0, errors.New("db")
		},
	}
	sErr := NewProductService(repoErr, client)
	_, _, err = sErr.Search(context.Background(), 1, 10, req)
	if err == nil {
		t.Fatalf("expected error")
	}
}
