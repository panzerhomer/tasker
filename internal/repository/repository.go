package repository

import (
	"context"
	"rest/internal/domain"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *domain.User) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, userID int64) (*domain.User, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, authorID int64, task *domain.Task) error
	DeleteTask(ctx context.Context, authorID int64, taskID int64) error
}

type ProjectRepository interface {
	CreateProject(ctx context.Context, authorID int64, project *domain.Project) (int64, error)
	JoinProjectByName(ctx context.Context, projectName string, userID int64, role string) error
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)
	GetProjectById(ctx context.Context, projectID int64) (*domain.Project, error)
}
