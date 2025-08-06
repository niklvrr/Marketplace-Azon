package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/pkg/models"
)

var (
	createProductQuery = `
		INSERT INTO products (seller_id, category_id, name, description, price, stock, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`

	getProductByIdQuery = `
		SELECT seller_id, category_id, name, description, price, stock, created_at
		FROM products 
		WHERE id = $1;`

	updateProductByIdQuery = `
		UPDATE products
		SET seller_id = $1, category_id = $2, name = $3, description = $4, price = $5, stock = $6, create_at = $7
		WHERE id = $8;`

	deleteProductByIdQuery = `DELETE FROM products WHERE id = $1;`

	getAllProductsQuery = `
		SELECT seller_id, category_id, name, description, price, stock, created_at
		FROM products
		ORDER BY name;`

	searchQuery = `
		SELECT seller_id, category_id, name, description, price, stock, created_at
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

func CreateProduct(ctx context.Context, db *pgxpool.Pool, p *models.Product) error {
	err := db.QueryRow(ctx, createProductQuery)
	if err != nil {
		return createProductError
	}

	return nil
}

func GetProductById(ctx context.Context, db *pgxpool.Pool, id int64) (*models.Product, error) {
	product := new(models.Product)
	err := db.QueryRow(ctx, getProductByIdQuery, id).
		Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt)

	if err != nil {
		return nil, productNotFound
	}

	return product, nil
}

func UpdateProductById(ctx context.Context, db *pgxpool.Pool, product *models.Product) error {
	cmdTag, err := db.Exec(
		ctx, updateProductByIdQuery,
		product.SellerId,
		product.CategoryId,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.Id)

	if err != nil {
		return updateProductError
	}

	if cmdTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}

func DeleteProductById(ctx context.Context, db *pgxpool.Pool, id int64) error {
	cmtTag, err := db.Exec(ctx, deleteProductByIdQuery, id)
	if err != nil {
		return deleteProductError
	}

	if cmtTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}

func GetAllProducts(ctx context.Context, db *pgxpool.Pool) (*[]models.Product, error) {
	rows, err := db.Query(ctx, getAllProductsQuery)
	if err != nil {
		return &[]models.Product{}, getAllProductsError
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt)

		if err != nil {
			return &[]models.Product{}, getAllProductsError
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]models.Product{}, rowsIterationError
	}

	return &products, nil
}

func SearchProducts(
	ctx context.Context, db *pgxpool.Pool,
	text string,
	categoryId *int64,
	min, max *float64,
	offset, limit int64,
) (*[]models.Product, error) {
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

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, searchProductsError
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err = rows.Scan(
			&product.SellerId,
			&product.CategoryId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
		)

		if err != nil {
			return &[]models.Product{}, searchProductsError
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return &[]models.Product{}, rowsIterationError
	}

	return &products, nil
}
