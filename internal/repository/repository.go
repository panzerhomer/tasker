package repository

import (
	"context"
	"rest/internal/domain"

	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	InsertUser(ctx context.Context, email string, password string) (int64, error)
	GetUser(ctx context.Context, email string) (*domain.User, error)
}

type Repository interface {
	UserRepository
}

func NewRepository(db *pgx.Conn) Repository {
	return NewUserRepository(db)
}
