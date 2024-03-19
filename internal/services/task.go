package services

import (
	"context"
	"fmt"
	"rest/internal/domain"
	"rest/internal/repository"
)

type taskService struct {
	taskRepo repository.TaskRepository
}

func NewTaskService(taskRepo repository.TaskRepository) *taskService {
	return &taskService{taskRepo: taskRepo}
}

func (s *taskService) CreateTask(ctx context.Context, authorID int64, task *domain.Task) error {
	err := s.taskRepo.CreateTask(ctx, authorID, task)
	if err != nil {
		return fmt.Errorf("%w: ", err)
	}

	return nil
}

// func (s *taskService) GetTask(ctx context.Context)

// func (s *taskService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
// 	user, err := s.userRepo.GetUserByEmail(ctx, email)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: ", err)
// 	}

// 	return user, nil
// }
