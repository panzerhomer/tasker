package repository

import (
	"context"
	"fmt"
	"rest/internal/domain"

	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	db *pgx.Conn
}

func NewUserRepo(db *pgx.Conn) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) InsertUser(ctx context.Context, user *domain.User) (int64, error) {
	const query = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id"

	stmt, err := r.db.Prepare(ctx, "insertUser", query)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	var id int64
	err = r.db.QueryRow(ctx, stmt.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	return id, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	const query = "SELECT user_id, email, password FROM users WHERE user_id = $1"

	var user domain.User

	stmt, err := r.db.Prepare(ctx, "selectUserById", query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = r.db.QueryRow(ctx, stmt.Name, userID).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &user, nil
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	const query = "SELECT user_id, email, password FROM users WHERE email = $1"

	var user domain.User

	stmt, err := r.db.Prepare(ctx, "selectUserByEmail", query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = r.db.QueryRow(ctx, stmt.Name, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &user, nil
}
