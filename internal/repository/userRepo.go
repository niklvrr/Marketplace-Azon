package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/internal/model"
)

var (
	createUserQuery = `
		INSERT INTO users (id, name, email, password, role, is_active, create_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, create_at`

	getUserByIdQuery = `
		SELECT id, name, email, password, role, is_active, create_at
		FROM users WHERE id = $1`

	getUserByEmailQuery = `
		SELECT id, name, email, password, role, is_active, create_at
		FROM users WHERE email = $1`

	blockUserByIdQuery = `UPDATE users SET is_active = $1 WHERE id = $2`
)

var (
	createUserError   = errors.New("error creating user")
	userNotFoundError = errors.New("user not found")
	blockExecError    = errors.New("error executing blocking user by id")
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
	).Scan(&user.Id, &user.CreateAt)

	if err != nil {
		return createUserError
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

func (r *UserRepo) BlockUserById(ctx context.Context, userId int64) error {
	cmdTag, err := r.db.Exec(ctx, blockUserByIdQuery, false, userId)
	if err != nil {
		return blockExecError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}
