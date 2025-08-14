package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createProductQuery = `
		INSERT INTO products (seller_id, category_id, name, description, price, stock, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id;`

	getProductByIdQuery = `
		SELECT seller_id, category_id, name, description, price, stock, status, created_at
		FROM products 
		WHERE id = $1;`

	updateProductByIdQuery = `
		UPDATE products
		SET seller_id = $1, category_id = $2, name = $3, description = $4, price = $5, stock = $6, status = $7, create_at = $8
		WHERE id = $9;`

	deleteProductByIdQuery = `DELETE FROM products WHERE id = $1;`

	countQuery = `SELECT COUNT(*) FROM products;`

	getAllProductsQuery = `
		SELECT seller_id, category_id, name, description, price, stock, status, created_at
		FROM products
		ORDER BY name
		LIMIT $1 OFFSET $2;`

	searchQuery = `
		SELECT seller_id, category_id, name, description, price, stock, status, created_at
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
	err := r.db.QueryRow(ctx, createProductQuery).Scan(&p.Id)
	if err != nil {
		return createProductError
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
		return nil, productNotFound
	}

	return product, nil
}

func (r *ProductRepo) UpdateProductById(ctx context.Context, product *model.Product) error {
	cmdTag, err := r.db.Exec(
		ctx, updateProductByIdQuery,
		product.SellerId,
		product.CategoryId,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.Status,
		product.Id)

	if err != nil {
		return updateProductError
	}

	if cmdTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}

func (r *ProductRepo) DeleteProductById(ctx context.Context, id int64) error {
	cmtTag, err := r.db.Exec(ctx, deleteProductByIdQuery, id)
	if err != nil {
		return deleteProductError
	}

	if cmtTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}

func (r *ProductRepo) GetAllProducts(ctx context.Context, offset, limit int) (*[]model.Product, int64, error) {
	rows, err := r.db.Query(ctx, getAllProductsQuery, limit, offset)
	if err != nil {
		return &[]model.Product{}, 0, getAllProductsError
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
			return &[]model.Product{}, 0, getAllProductsError
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]model.Product{}, 0, rowsIterationError
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return &[]model.Product{}, 0, getAllProductsError
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
		return nil, 0, searchProductsError
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
			return &[]model.Product{}, 0, searchProductsError
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]model.Product{}, 0, rowsIterationError
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
		return &[]model.Product{}, 0, searchProductsError
	}

	return &products, total, nil
}
