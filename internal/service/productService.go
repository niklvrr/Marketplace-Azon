package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/niklvrr/myMarketplace/internal/model"
	"github.com/redis/go-redis/v9"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *model.Product) error
	GetProductById(ctx context.Context, productId int64) (*model.Product, error)
	UpdateProductById(ctx context.Context, product *model.Product) error
	DeleteProductById(ctx context.Context, productId int64) error
	GetAllProducts(ctx context.Context, offset, limit int) (*[]model.Product, int64, error)
	SearchProducts(
		ctx context.Context,
		text *string,
		categoryId *int64,
		min, max *float64,
		offset, limit int,
	) (*[]model.Product, int64, error)
}

type ProductService struct {
	repo  IProductRepository
	cache *redis.Client
}

func NewProductService(repo IProductRepository, cache *redis.Client) *ProductService {
	return &ProductService{
		repo:  repo,
		cache: cache,
	}
}

func (s *ProductService) Create(ctx context.Context, sellerId int64, req *model.CreateProductRequest) (model.ProductResponse, error) {
	p := model.Product{
		SellerId:    sellerId,
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	err := s.repo.CreateProduct(ctx, &p)
	if err != nil {
		return model.ProductResponse{}, err
	}

	s.cache.Del(ctx, "products:all")

	return model.ProductResponse{
		Id:          p.Id,
		SellerId:    p.SellerId,
		CategoryId:  p.CategoryId,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}, nil
}

func (s *ProductService) GetById(ctx context.Context, req *model.GetProductsRequest) (model.ProductResponse, error) {
	resp, err := s.repo.GetProductById(ctx, req.Id)
	if err != nil {
		return model.ProductResponse{}, err
	}

	return model.ProductResponse{
		Id:          resp.Id,
		SellerId:    resp.SellerId,
		CategoryId:  resp.CategoryId,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Stock:       resp.Stock,
	}, nil
}

func (s *ProductService) UpdateById(ctx context.Context, sellerId int64, req *model.UpdateProductRequest) (model.ProductResponse, error) {
	p := model.Product{
		SellerId:    sellerId,
		CategoryId:  *req.CategoryId,
		Name:        *req.Name,
		Description: *req.Description,
		Price:       *req.Price,
		Stock:       *req.Stock,
	}

	err := s.repo.UpdateProductById(ctx, &p)
	if err != nil {
		return model.ProductResponse{}, err
	}

	s.cache.Del(ctx, "products:all")
	return model.ProductResponse{
		Id:          p.Id,
		SellerId:    p.SellerId,
		CategoryId:  p.CategoryId,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}, nil
}

func (s *ProductService) DeleteById(ctx context.Context, req *model.DeleteProductRequest) error {
	err := s.repo.DeleteProductById(ctx, req.Id)
	if err != nil {
		return err
	}

	s.cache.Del(ctx, "products:all")
	return nil
}

func (s *ProductService) GetAll(ctx context.Context, page, limit int) ([]model.ProductResponse, int64, error) {
	const cachedKey = "products:all"

	cachedData, err := s.cache.Get(ctx, cachedKey).Bytes()
	if err == nil {
		var products []model.ProductResponse
		err = json.Unmarshal(cachedData, &products)
		if err == nil {
			return products, int64(len(products)), nil
		}
	}

	offset := (page - 1) * limit
	products, total, err := s.repo.GetAllProducts(ctx, offset, limit)
	if err != nil {
		return []model.ProductResponse{}, 0, err
	}

	var result []model.ProductResponse
	for _, product := range *products {
		resp := model.ProductResponse{
			Id:          product.Id,
			SellerId:    product.SellerId,
			CategoryId:  product.CategoryId,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		}

		result = append(result, resp)
	}

	dataToCache, err := json.Marshal(result)
	if err != nil {
		return []model.ProductResponse{}, 0, err
	}

	cacheExpiration := 5 * time.Minute
	s.cache.Set(ctx, cachedKey, string(dataToCache), cacheExpiration)

	return result, total, nil
}

func (s *ProductService) Search(ctx context.Context, page, limit int, req *model.SearchProductsRequest) ([]model.ProductResponse, int64, error) {
	offset := (page - 1) * limit
	products, total, err := s.repo.SearchProducts(
		ctx,
		req.Text,
		req.CategoryId,
		req.Min, req.Max,
		offset, limit,
	)
	if err != nil {
		return []model.ProductResponse{}, 0, err
	}

	var result []model.ProductResponse
	for _, product := range *products {
		resp := model.ProductResponse{
			Id:          product.Id,
			SellerId:    product.SellerId,
			CategoryId:  product.CategoryId,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Stock:       product.Stock,
		}

		result = append(result, resp)
	}

	return result, total, nil
}
