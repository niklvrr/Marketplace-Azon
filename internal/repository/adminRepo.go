package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	getAllUsersQuery = `
		SELECT user_id, name, email, password, role, is_active, create_at
		FROM users`

	updateUserRoleQuery = `UPDATE users SET role=$1 WHERE user_id=$2`

	updateUserStatusQuery = `UPDATE users SET is_active=$1 WHERE user_id=$2`

	approveProductQuery = `UPDATE products SET is_approve=$1 WHERE product_id=$2`
)

var (
	getAllUsersError      = errors.New("GetAll Users Error")
	updateUserRoleError   = errors.New("update User Role Error")
	updateUserStatusError = errors.New("update User Status Error")
	approveProductError   = errors.New("approve Product Error")
)

type AdminRepo struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{db: db}
}

func (r *AdminRepo) GetAllUsers(ctx context.Context) ([]model.User, error) {
	rows, err := r.db.Query(ctx, getAllUsersQuery)
	if err != nil {
		return []model.User{}, getAllUsersError
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsActive,
			&user.CreateAt,
		)

		if err != nil {
			return []model.User{}, getAllUsersError
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return []model.User{}, rowsIterationError
	}

	return users, nil
}

func (r *AdminRepo) UpdateUserRole(ctx context.Context, userId int64, newRole string) error {
	cmdTag, err := r.db.Exec(ctx, updateUserRoleQuery, newRole, userId)
	if err != nil {
		return updateUserRoleError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}

func (r *AdminRepo) UpdateUserStatus(ctx context.Context, userId int64, newStatus bool) error {
	cmdTag, err := r.db.Exec(ctx, updateUserStatusQuery, newStatus, userId)
	if err != nil {
		return updateUserStatusError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}

func (r *AdminRepo) ApproveProduct(ctx context.Context, userId int64) error {
	cmdTag, err := r.db.Exec(ctx, approveProductQuery, true, userId)
	if err != nil {
		return approveProductError
	}

	if cmdTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}
