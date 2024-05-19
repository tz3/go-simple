package repository

//go:generate mockgen -source user.go -package repository -destination user_mock.go

import (
	"context"
	"database/sql"
	"errors"

	"github.com/tz3/go-simple/user/model"
)

var (
	// ErrUserNotFound is returned when a user is not found in the database
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int) (*model.User, error)
}

type UserRepositoryImpl struct {
	DB *sql.DB
}

// NewUserRepository will return a new userRepo.
func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepositoryImpl{
		DB: db,
	}
}

// CreateUser will create a new user in the database and rewrite it the user argument.
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (first_name, last_name) VALUES ($1, $2) RETURNING id"

	err := r.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID will return a new user from the database by a given ID.
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	query := "SELECT first_name, last_name, id FROM users WHERE id = $1"

	row := r.DB.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.FirstName, &user.LastName, &user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
