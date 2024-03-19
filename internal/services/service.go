package services

import (
	"context"
	"rest/internal/domain"
)

type TaskService interface {
	CreateTask(ctx context.Context, UserID int64, task *domain.Task) error
	// GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type UserService interface {
	Login(ctx context.Context, user *domain.User) (string, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}
