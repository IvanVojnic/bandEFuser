package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"userMS/models"
)

// UserPostgres has an internal db object
type UserPostgres struct {
	db *pgxpool.Pool
}

// NewUserAuthPostgres used to init UsesAP
func NewUserPostgres(db *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{db: db}
}

// SignUp used to create user
func (r *UserPostgres) SignUp(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, "insert into users (id, email, password) values($1, $2, $3)",
		user.UserID, user.UserEmail, user.Password)
	if err != nil {
		return fmt.Errorf("error while user creating: %v", err)
	}
	return nil
}

// UpdateRefreshToken used to update rt
func (r *UserPostgres) UpdateRefreshToken(ctx context.Context, rt string, id uuid.UUID) error {
	_, errInsert := r.db.Exec(ctx, "UPDATE users SET refreshtoken = $1 WHERE id = $2", rt, id)
	if errInsert != nil {
		return fmt.Errorf("update user error %w", errInsert)
	}
	return nil
}

// GetUserByID used to get user by ID
func (r *UserPostgres) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	user := models.User{}
	err := r.db.QueryRow(ctx,
		"select users.id, users.email, users.password, users.refreshtoken from users where id=$1",
		userID).Scan(&user.UserID, &user.UserEmail, &user.Password, &user.RefreshToken)
	if err != nil {
		return user, fmt.Errorf("get user error %w", err)
	}

	return user, nil
}

// SignInUser used to sign in user
func (r *UserPostgres) SignIn(ctx context.Context, user *models.User) error {
	err := r.db.QueryRow(ctx,
		`select users.id, users.email, users.password from users where email=$1`,
		user.UserEmail).Scan(&user.UserID, &user.UserEmail, &user.Password)
	if err != nil {
		return fmt.Errorf("error while getting user %w", err)
	}
	return nil
}
