package services

import (
	"context"
	"rest/internal/domain"
	"rest/internal/lib/hasher"
	"rest/internal/lib/jwt"
	"rest/internal/repository"
	"time"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) Login(ctx context.Context, user *domain.User) (string, error) {
	u, err := us.repo.GetUser(ctx, user.Email)
	if err != nil {
		return "", err
	}
	err = hasher.CheckPassword(user.Password, u.Password)
	if err != nil {
		return "", err
	}

	token, err := jwt.NewToken(u, time.Minute)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (us *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// _, err := us.repo.GetUser(ctx, user.Email)
	// if err == nil {
	// 	return nil, fmt.Errorf("user exist")
	// }

	hashedPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	u := &domain.User{
		Email:    user.Email,
		Password: hashedPassword,
	}

	id, err := us.repo.InsertUser(ctx, u.Email, u.Password)
	if err != nil {
		return nil, err
	}

	u.ID = id

	return u, nil
}
