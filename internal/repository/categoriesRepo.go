package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createCategoryQuery = `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id;`

	getCategoryByIdQuery = `SELECT id, name, description FROM categories WHERE id = $1;`

	updateCategoryByIdQuery = `UPDATE categories SET name = $1, description = $2 WHERE id = $3;`

	deleteCategoryByIdQuery = `DELETE FROM categories WHERE id = $1;`

	getAllCategoriesQuery = `
		SELECT id, name, description 
		FROM categories
		ORDER BY name`
)

var (
	createCategoryError   = errors.New("create category error")
	categoryNotFoundError = errors.New("category not found")
	updateCategoryError   = errors.New("update category error")
	deleteCategoryError   = errors.New("delete category error")
	getAllCategoriesError = errors.New("get all categories error")
	rowsIterationError    = errors.New("rows iteration error")
)

type CategoryRepo struct {
	db *pgxpool.Pool
}

func NewCategoryRepo(db *pgxpool.Pool) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) CreateCategory(ctx context.Context, category *model.Category) error {
	err := r.db.QueryRow(
		ctx, createCategoryQuery, category.Name, category.Description).
		Scan(&category.Id)

	if err != nil {
		return fmt.Errorf("%w: %w", createCategoryError, err)
	}

	return nil
}

func (r *CategoryRepo) GetCategoryById(ctx context.Context, categoryId int64) (*model.Category, error) {
	category := new(model.Category)
	err := r.db.QueryRow(ctx, getCategoryByIdQuery, categoryId).Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return &model.Category{}, fmt.Errorf("%w: %w", categoryNotFoundError, err)
	}

	return category, nil
}

func (r *CategoryRepo) UpdateCategory(ctx context.Context, category model.Category) error {
	cmdTag, err := r.db.Exec(ctx, updateCategoryByIdQuery, category.Name, category.Description, category.Id)
	if err != nil {
		return fmt.Errorf("%w: %w", updateCategoryError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", updateCategoryError, categoryNotFoundError)
	}

	return nil
}

func (r *CategoryRepo) DeleteCategory(ctx context.Context, categoryId int64) error {
	cmdTag, err := r.db.Exec(ctx, deleteCategoryByIdQuery, categoryId)
	if err != nil {
		return fmt.Errorf("%w: %w", deleteCategoryError, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", deleteCategoryError, categoryNotFoundError)
	}

	return nil
}

func (r *CategoryRepo) GetAllCategories(ctx context.Context) (*[]model.Category, error) {
	rows, err := r.db.Query(ctx, getAllCategoriesQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", getAllCategoriesError, err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		category := new(model.Category)
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", getAllCategoriesError, err)
		}

		categories = append(categories, *category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w(%w): %w", getAllCategoriesError, rowsIterationError, err)
	}

	return &categories, nil
}
