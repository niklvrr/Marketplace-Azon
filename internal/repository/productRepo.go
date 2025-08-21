package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createProductQuery = `
		INSERT INTO products (seller_id, category_id, name, description, price, stock, is_approved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`

	getProductByIdQuery = `
		SELECT seller_id, category_id, name, description, price, stock, is_approved, created_at
		FROM products 
		WHERE id = $1;`

	updateProductByIdQuery = `
		UPDATE products
		SET category_id = $1, name = $2, description = $3, price = $4, stock = $5,
		WHERE id = $6;`

	deleteProductByIdQuery = `DELETE FROM products WHERE id = $1;`

	countQuery = `SELECT COUNT(*) FROM products;`

	getAllProductsQuery = `
		SELECT seller_id, category_id, name, description, price, stock, is_approved, created_at
		FROM products
		ORDER BY name
		LIMIT $1 OFFSET $2;`

	searchQuery = `
		SELECT seller_id, category_id, name, description, price, stock, is_approved, created_at
		FROM products
		WHERE to_tsvector('simple', name || ' ' || coalesce(description, '')) @@ plainto_tsquery('simple', $1)`
)

var (
	createProductError  = errors.New(`error creating product`)
	productNotFound     = errors.New(`product not found`)
	updateProductError  = errors.New(`error updating product`)
	deleteProductError  = errors.New(`error deleting product`)
	getAllProductsError = errors.New(`error getting all products`)
	searchProductsError = errors.New(`error searching products`)
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) CreateProduct(ctx context.Context, p *model.Product) error {
	err := r.db.QueryRow(
		ctx, createProductQuery,
		p.SellerId,
		p.CategoryId,
		p.Name,
		p.Description,
		p.Price,
		p.Stock,
		false,
		time.Now(),
	).Scan(&p.Id)
	if err != nil {
		return fmt.Errorf("%w: %w", createProductError, err)
	}

	return nil
}

func (r *ProductRepo) GetProductById(ctx context.Context, id int64) (*model.Product, error) {
	product := new(model.Product)
	err := r.db.QueryRow(ctx, getProductByIdQuery, id).
		Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Status,
			&product.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", productNotFound, err)
	}

	return product, nil
}

func (r *ProductRepo) UpdateProductById(ctx context.Context, product *model.Product) error {
	cmdTag, err := r.db.Exec(
		ctx, updateProductByIdQuery,
		product.CategoryId,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Id)

	if err != nil {
		return fmt.Errorf("%w: %w", updateProductError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", updateProductError, productNotFound)
	}

	return nil
}

func (r *ProductRepo) DeleteProductById(ctx context.Context, id int64) error {
	cmtTag, err := r.db.Exec(ctx, deleteProductByIdQuery, id)
	if err != nil {
		return deleteProductError
	}

	if cmtTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", deleteProductError, productNotFound)
	}

	return nil
}

func (r *ProductRepo) GetAllProducts(ctx context.Context, offset, limit int) (*[]model.Product, int64, error) {
	rows, err := r.db.Query(ctx, getAllProductsQuery, limit, offset)
	if err != nil {
		return &[]model.Product{}, 0, fmt.Errorf("%w: %w", getAllProductsError, err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err = rows.Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Status,
			&product.CreatedAt)

		if err != nil {
			return &[]model.Product{}, 0, fmt.Errorf("%w: %w", getAllProductsError, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]model.Product{}, 0, fmt.Errorf("%w(%w): %w", getAllProductsError, rowsIterationError, err)
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return &[]model.Product{}, 0, fmt.Errorf("%w: %w", getAllProductsError, err)
	}

	return &products, total, nil
}

func (r *ProductRepo) SearchProducts(
	ctx context.Context,
	text *string,
	categoryId *int64,
	min, max *float64,
	offset, limit int,
) (*[]model.Product, int64, error) {
	sql := searchQuery

	args := []interface{}{text}

	if categoryId != nil {
		sql += fmt.Sprintf(" AND category_id = $%d", len(args))
	}

	if min != nil {
		sql += fmt.Sprintf(" AND min = $%d", len(args))
	}

	if max != nil {
		sql += fmt.Sprintf(" AND max = $%d", len(args))
	}

	args = append(args, offset, limit)
	sql += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %w", searchProductsError, err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		err = rows.Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.Status,
			&product.CreatedAt,
		)

		if err != nil {
			return &[]model.Product{}, 0, fmt.Errorf("%w: %w", searchProductsError, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]model.Product{}, 0, fmt.Errorf("%w(%w): %w", searchProductsError, rowsIterationError, err)
	}

	var where []string
	var countArgs []interface{}

	param := func() string { return fmt.Sprintf("$%d", len(countArgs)+1) }

	if text != nil && strings.TrimSpace(*text) != "" {
		where = append(where, "title ILIKE '%' || "+param()+" || '%'")
		countArgs = append(countArgs, *text)
	}

	if categoryId != nil && *categoryId > 0 {
		where = append(where, "category_id = "+param())
		countArgs = append(countArgs, *categoryId)
	}

	if min != nil {
		where = append(where, "price >= "+param())
		countArgs = append(countArgs, *min)
	}

	if max != nil {
		where = append(where, "price <= "+param())
		countArgs = append(countArgs, *max)
	}

	var q string
	if len(where) > 0 {
		q = countQuery + " WHERE " + strings.Join(where, " AND ")
	} else {
		q = countQuery
	}

	var total int64
	err = r.db.QueryRow(ctx, q, countArgs...).Scan(&total)
	if err != nil {
		return &[]model.Product{}, 0, fmt.Errorf("%w: %w", searchProductsError, err)
	}

	return &products, total, nil
}
