package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niklvrr/myMarketplace/pkg/models"
)

var (
	createUserQuery = `
		INSERT INTO users (id, name, email, password, role, create_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, role, create_at`

	getUserByIdQuery = `
		SELECT id, name, email, password, role, create_at
		FROM users WHERE id = $1`

	getUserByEmailQuery = `
		SELECT id, name, email, password, role, create_at
		FROM users WHERE email = $1`

	blockUserByIdQuery = `UPDATE users SET role = $1 WHERE id = $2`
)

var (
	createUserError   = errors.New("error creating user")
	userNotFoundError = errors.New("user not found")
	blockExecError    = errors.New("error executing blocking user by id")
)

func CreateUser(ctx context.Context, db *pgxpool.Pool, user *models.User) error {
	err := db.QueryRow(
		ctx, createUserQuery,
	).Scan(&user.Id, &user.CreateAt)

	if err != nil {
		return createUserError
	}

	return nil
}

func GetUserById(ctx context.Context, db *pgxpool.Pool, userId int64) (*models.User, error) {
	user := new(models.User)
	err := db.QueryRow(ctx, getUserByIdQuery, userId).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt)

	if err != nil {
		return &models.User{}, userNotFoundError
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*models.User, error) {
	user := new(models.User)
	err := db.QueryRow(ctx, getUserByEmailQuery, email).
		Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreateAt)

	if err != nil {
		return &models.User{}, userNotFoundError
	}

	return user, nil
}

func BlockUserById(ctx context.Context, db *pgxpool.Pool, userId int64) error {
	cmdTag, err := db.Exec(ctx, blockUserByIdQuery, "blocked", userId)
	if err != nil {
		return blockExecError
	}

	if cmdTag.RowsAffected() == 0 {
		return userNotFoundError
	}

	return nil
}
