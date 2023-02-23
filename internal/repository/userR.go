package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"userMS/models"
)

// UserCommPostgres is a wrapper to db object
type UserCommPostgres struct {
	db *pgxpool.Pool
}

type Status int

const (
	Decline Status = iota
	Accept
	NoAnswer
)

// NewUserCommPostgres used to init UserCommPostgres
func NewUserCommPostgres(db *pgxpool.Pool) *UserCommPostgres {
	return &UserCommPostgres{db: db}
}

// GetFriends used to send friends
func (r *UserCommPostgres) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	users := make([]models.User, 0)
	rowsFriends, err := r.db.Query(ctx,
		`select users.id, users.name from users u
    		inner join friends f on f.userReceiver = u.id or f.userSender = u.id where u.id=$1 and f.status=$2`, userID, Accept)
	if err != nil {
		return users, fmt.Errorf("get all friends sql script error %w", err)
	}
	defer rowsFriends.Close()
	for rowsFriends.Next() {
		var user models.User
		errScan := rowsFriends.Scan(&user.ID, &user.Name)
		if errScan != nil {
			return users, fmt.Errorf("get all friends scan rows error %w", errScan)
		}
		users = append(users, user)
	}
	return users, nil
}

// SendFriendsRequest used to send requests for user
func (r *UserCommPostgres) SendFriendsRequest(ctx context.Context, userSender uuid.UUID, userReceiver uuid.UUID) error {
	friendsID := uuid.New()
	_, err := r.db.Exec(ctx, "insert into friends (id, userSender, userReceiver, status_id) values($1, $2, $3, $4)",
		friendsID, userSender, userReceiver, NoAnswer)
	if err != nil {
		return fmt.Errorf("error while friends relationship creating: %s", err)
	}
	return nil
}

// AcceptFriendsRequest used to accept request
func (r *UserCommPostgres) AcceptFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE friends 
			SET status=$1 
			WHERE userSender=$2 AND userReceiver=$3`,
		Accept, userSenderID, userID)
	if err != nil {
		return fmt.Errorf("update friends error %w", err)
	}
	return nil
}

// AcceptFriendsRequest used to accept request
func (r *UserCommPostgres) DeclineFriendsRequest(ctx context.Context, userSenderID uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE friends 
			SET status=$1 
			WHERE userSender=$2 AND userReceiver=$3`,
		Decline, userSenderID, userID)
	if err != nil {
		return fmt.Errorf("update friends error %w", err)
	}
	return nil
}

// FindUser used to find user by email
func (r *UserCommPostgres) FindUser(ctx context.Context, userEmail string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx, "select users.id, users.email from users where users.email=$1", userEmail).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, fmt.Errorf("error: cannot get id, %w", err)
	}
	return user, nil
}

// GetRequest used to send request to be a friends
func (r *UserCommPostgres) GetRequest(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	users := make([]models.User, 0)
	rowsFriendsReq, err := r.db.Query(ctx,
		`select users.id, users.name from users u
    		inner join friends f on f.userReceiver = u.id where u.id=$1 and f.status=$2`, userID, NoAnswer)
	if err != nil {
		return users, fmt.Errorf("get all friends requests sql script error %w", err)
	}
	defer rowsFriendsReq.Close()
	for rowsFriendsReq.Next() {
		var user models.User
		errScan := rowsFriendsReq.Scan(&user.ID, &user.Name)
		if errScan != nil {
			return users, fmt.Errorf("get all friends requests scan rows error %w", errScan)
		}
		users = append(users, user)
	}
	return users, nil
}
