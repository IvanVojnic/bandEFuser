package repository

import (
	"context"
	"fmt"
	"github.com/IvanVojnic/bandEFuser/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	_, err := r.db.Exec(ctx, "insert into users (id, email, name, password) values($1, $2, $3, $4)",
		user.ID, user.Email, user.Name, user.Password)
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
		"select users.id, users.email, users.name, users.password, users.refreshtoken from users where users.id=$1",
		userID).Scan(&user.ID, &user.Email, &user.Name, &user.Password, &user.RefreshToken)
	if err != nil {
		return user, fmt.Errorf("get user error %w", err)
	}

	return user, nil
}

// SignInUser used to sign in user
func (r *UserPostgres) SignIn(ctx context.Context, user *models.User) error {
	err := r.db.QueryRow(ctx,
		`select users.id, users.name, users.password from users where users.name=$1`,
		user.Name).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		return fmt.Errorf("error while getting user %w", err)
	}
	return nil
}
