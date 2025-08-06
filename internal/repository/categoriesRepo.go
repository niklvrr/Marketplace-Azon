package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/pkg/models"
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

func CreateCategory(ctx context.Context, db *pgxpool.Pool, category *models.Category) error {
	err := db.QueryRow(
		ctx, createCategoryQuery, category.Name, category.Description).
		Scan(&category.Id)

	if err != nil {
		return createCategoryError
	}

	return nil
}

func GetCategoryById(ctx context.Context, db *pgxpool.Pool, categoryId int64) (*models.Category, error) {
	category := new(models.Category)
	err := db.QueryRow(ctx, getCategoryByIdQuery, categoryId).Scan(&category.Id, &category.Name, &category.Description)
	if err != nil {
		return &models.Category{}, categoryNotFoundError
	}

	return category, nil
}

func UpdateCategory(ctx context.Context, db *pgxpool.Pool, category models.Category) error {
	cmdTag, err := db.Exec(ctx, updateCategoryByIdQuery, category.Name, category.Description, category.Id)
	if err != nil {
		return updateCategoryError
	}

	if cmdTag.RowsAffected() == 0 {
		return categoryNotFoundError
	}

	return nil
}

func DeleteCategory(ctx context.Context, db *pgxpool.Pool, categoryId int64) error {
	cmdTag, err := db.Exec(ctx, deleteCategoryByIdQuery, categoryId)
	if err != nil {
		return deleteCategoryError
	}

	if cmdTag.RowsAffected() == 0 {
		return categoryNotFoundError
	}

	return nil
}

func GetAllCategories(ctx context.Context, db *pgxpool.Pool) (*[]models.Category, error) {
	rows, err := db.Query(ctx, getAllCategoriesQuery)
	if err != nil {
		return nil, getAllCategoriesError
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		category := new(models.Category)
		err := rows.Scan(&category.Id, &category.Name, &category.Description)
		if err != nil {
			return nil, getAllCategoriesError
		}

		categories = append(categories, *category)
	}

	if err := rows.Err(); err != nil {
		return nil, rowsIterationError
	}

	return &categories, nil
}
