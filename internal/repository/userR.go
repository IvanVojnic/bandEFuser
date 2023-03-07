// Package repository consists of communicates methods for user
package repository

import (
	"context"
	"fmt"

	"github.com/IvanVojnic/bandEFuser/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserCommPostgres is a wrapper to db object
type UserCommPostgres struct {
	db *pgxpool.Pool
}

// Status used to define types of statuses
type Status int

// Decline define negative status
// Accept define positive status
// NoAnswer define neutral status
const (
	Decline Status = iota
	Accept
	NoAnswer
)

// NewUserCommPostgres used to init UserCommPostgres
func NewUserCommPostgres(db *pgxpool.Pool) *UserCommPostgres {
	return &UserCommPostgres{db: db}
}

// GetFriends used to get friends
func (r *UserCommPostgres) GetFriends(ctx context.Context, userID uuid.UUID) (*[]models.User, error) { // nolint:dupl, gocritic
	var users *[]models.User
	rowsFriends, err := r.db.Query(ctx,
		`SELECT users.id, users.name FROM users
    		INNER JOIN friends ON friends.userReceiver = users.id OR friends.userSender = users.id WHERE users.id=$1 AND friends.status=$2`, userID, Accept)
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
		*users = append(*users, user)
	}
	return users, nil
}

// SendFriendsRequest used to send requests for user
func (r *UserCommPostgres) SendFriendsRequest(ctx context.Context, userSender, userReceiver uuid.UUID) error {
	friendsID := uuid.New()
	_, err := r.db.Exec(ctx, `INSERT INTO friends (id, userSender, userReceiver, status_id) VALUES($1, $2, $3, $4)`,
		friendsID, userSender, userReceiver, NoAnswer)
	if err != nil {
		return fmt.Errorf("error while friends relationship creating: %s", err)
	}
	return nil
}

// AcceptFriendsRequest used to accept request
func (r *UserCommPostgres) AcceptFriendsRequest(ctx context.Context, userSenderID, userID uuid.UUID) error {
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

// DeclineFriendsRequest used to decline request
func (r *UserCommPostgres) DeclineFriendsRequest(ctx context.Context, userSenderID, userID uuid.UUID) error {
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
func (r *UserCommPostgres) FindUser(ctx context.Context, userEmail string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx, `SELECT users.id, users.email FROM users WHERE users.email=$1`, userEmail).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return &user, fmt.Errorf("error: cannot get id, %w", err)
	}
	return &user, nil
}

// GetRequest used to get request to be a friends
func (r *UserCommPostgres) GetRequest(ctx context.Context, userID uuid.UUID) (*[]models.User, error) { // nolint:dupl, gocritic
	var users *[]models.User
	rowsFriendsReq, err := r.db.Query(ctx,
		`SELECT users.id, users.name FROM users u
    		INNER JOIN friends f ON f.userReceiver = u.id WHERE u.id=$1 AND f.status=$2`, userID, NoAnswer)
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
		*users = append(*users, user)
	}
	return users, nil
}

func (r *UserCommPostgres) GetUsers(ctx context.Context, usersID *[]uuid.UUID) (*[]models.User, error) {
	var users []models.User
	rowsUsersReq, err := r.db.Query(ctx,
		`SELECT users.id, users.name, users.email 
			 FROM users 
			 WHERE users.id = ANY($1::uuid[])`, usersID)
	if err != nil {
		return nil, fmt.Errorf("error while getting users from db, %s", err)
	}
	defer rowsUsersReq.Close()
	for rowsUsersReq.Next() {
		var user models.User
		errScan := rowsUsersReq.Scan(&user.ID, &user.Name, &user.Email)
		if errScan != nil {
			return nil, fmt.Errorf("get all friends requests scan rows error %w", errScan)
		}
		users = append(users, user)
	}
	return &users, nil
}
