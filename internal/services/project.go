package services

import (
	"context"
	"errors"
	"log"
	"rest/internal/domain"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, authorID int64, project domain.Project) (int64, error)
	JoinProjectByName(ctx context.Context, projectName string, userID int64, role int8) error
	GetProjectByName(ctx context.Context, name string) (*domain.Project, error)
	GetProjectById(ctx context.Context, projectID int64) (*domain.Project, error)
	GetProjects(ctx context.Context, userID int64) ([]*domain.Project, error)
}

type projectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *projectService {
	return &projectService{repo: repo}
}

func (s *projectService) CreateProject(ctx context.Context, authorID int64, project domain.Project) error {
	_, err := s.repo.CreateProject(ctx, authorID, project)
	if err != nil {
		return errors.New("can't create project")
	}

	return nil
}

func (s *projectService) JoinProject(ctx context.Context, projectName string, userID int64, role int8) error {
	err := s.repo.JoinProjectByName(ctx, projectName, userID, role)
	log.Println("[serive get all projects]", err)
	if err != nil {
		return errors.New("can't join project")
	}

	return nil
}

func (s *projectService) GetAllProjects(ctx context.Context, userID int64) ([]*domain.Project, error) {
	projects, err := s.repo.GetProjects(ctx, userID)
	log.Println("[serive get all projects]", err)
	if err != nil {
		return nil, errors.New("can't get all project")
	}

	return projects, nil
}
