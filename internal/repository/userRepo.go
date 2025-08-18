package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createUserQuery = `
		INSERT INTO users (name, email, password, role, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	getUserByIdQuery = `
		SELECT id, name, email, password, role, is_active, created_at
		FROM users WHERE id = $1`

	getUserByEmailQuery = `
		SELECT id, name, email, password, role, is_active, created_at
		FROM users WHERE email = $1`

	updateUserByIdQuery = `
		UPDATE users
		SET name = $1, email = $2, password = $3
		WHERE id = $4`

	blockUserByIdQuery = `UPDATE users SET is_active = FALSE WHERE id = $1`

	unBlockUserByIdQuery = `UPDATE users SET is_active = TRUE WHERE id = $1`

	getAllUsersQuery = `
		SELECT id, name, email, password, role, is_active, created_at
		FROM users`

	updateUserRoleQuery = `UPDATE users SET role=$1 WHERE id=$2`

	approveProductQuery = `UPDATE products SET is_approve=TRUE WHERE product_id=$2`
)

var (
	createUserError     = errors.New("error creating user я того рот ебал")
	userNotFoundError   = errors.New("user not found")
	updateUserError     = errors.New("error updating user")
	blockExecError      = errors.New("error executing blocking user by id")
	unBlockExecError    = errors.New("error executing unblocking user by id")
	getAllUsersError    = errors.New("GetAll Users Error")
	updateUserRoleError = errors.New("update User Role Error")
	approveProductError = errors.New("approve Product Error")
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	err := r.db.QueryRow(
		ctx, createUserQuery,
		user.Name,
		user.Email,
		user.Password,
		"user",
		true,
		time.Now(),
	).Scan(&user.Id)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	user := new(model.User)
	err := r.db.QueryRow(ctx, getUserByIdQuery, userId).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsActive,
			&user.CreateAt)

	if err != nil {
		return &model.User{}, userNotFoundError
	}

	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.QueryRow(ctx, getUserByEmailQuery, email).
		Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.IsActive,
			&user.CreateAt)

	if err != nil {
		return &model.User{}, userNotFoundError
	}

	return user, nil
}

func (r *UserRepo) UpdateUserById(ctx context.Context, user *model.User) error {
	cmdTag, err := r.db.Exec(ctx, updateUserByIdQuery, user.Name, user.Email, user.Password, user.Id)
	if err != nil {
		return updateUserError
	}

	if cmdTag.RowsAffected() == 0 {
		return rowsIterationError
	}

	return nil
}

func (r *UserRepo) BlockUserById(ctx context.Context, userId int64) error {
	cmdTag, err := r.db.Exec(ctx, blockUserByIdQuery, userId)
	if err != nil {
		return blockExecError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}

func (r *UserRepo) UnBlockUserById(ctx context.Context, userId int64) error {
	cmdTag, err := r.db.Exec(ctx, unBlockUserByIdQuery, userId)
	if err != nil {
		return unBlockExecError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]model.User, error) {
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

func (r *UserRepo) UpdateUserRole(ctx context.Context, userId int64, newRole string) error {
	cmdTag, err := r.db.Exec(ctx, updateUserRoleQuery, newRole, userId)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}

func (r *UserRepo) ApproveProduct(ctx context.Context, productId int64) error {
	cmdTag, err := r.db.Exec(ctx, approveProductQuery, productId)
	if err != nil {
		return approveProductError
	}

	if cmdTag.RowsAffected() == 0 {
		return productNotFound
	}

	return nil
}
