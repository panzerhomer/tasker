package services

import (
	"context"
	"rest/internal/domain"
	"rest/internal/lib/hasher"
	"rest/internal/lib/jwt"
	"rest/internal/repository"
	"time"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) Login(ctx context.Context, user *domain.User) (string, error) {
	u, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}

	err = hasher.CheckPassword(user.Password, u.Password)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewToken(u, time.Minute*15)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Email:    user.Email,
		Password: hashedPassword,
	}

	id, err := s.repo.InsertUser(ctx, u)
	if err != nil {
		return nil, err
	}

	u.ID = id

	return u, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
