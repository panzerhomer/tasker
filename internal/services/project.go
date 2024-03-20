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
	LeaveProjectByName(ctx context.Context, projectName string, userID int64) error
	IsMember(ctx context.Context, userID int64, projectName string) (int8, error) // return role and error
	GetProjectUsers(ctx context.Context, userID int64, projectName string) ([]*domain.ProjectUser, error)
}

type TaskRepository interface {
	AssignUserToTask(ctx context.Context, userID int64, projectID int64, task domain.Task) error
	UpdateTask(ctx context.Context, authorID int64, taskID int64) error
}

type projectService struct {
	repo     ProjectRepository
	taksRepo TaskRepository
}

func NewProjectService(repo ProjectRepository, taksRepo TaskRepository) *projectService {
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

func (s *projectService) LeaveProject(ctx context.Context, projectName string, userID int64) error {
	err := s.repo.LeaveProjectByName(ctx, projectName, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *projectService) GetAllProjects(ctx context.Context, userID int64) ([]*domain.Project, error) {
	projects, err := s.repo.GetProjects(ctx, userID)
	if err != nil {
		return nil, errors.New("can't get all project")
	}

	return projects, nil
}

func (s *projectService) CreateProjectTask(ctx context.Context, userID int64, projectName string, task domain.Task) error {
	role, err := s.repo.IsMember(ctx, userID, projectName)
	if err != nil {
		return err
	}

	if role != domain.AdminRole {
		return errors.New("must be admin to create task")
	}

	project, err := s.repo.GetProjectByName(ctx, projectName)
	if err != nil {
		return err
	}

	err = s.taksRepo.AssignUserToTask(ctx, userID, project.ID, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *projectService) GetProjectUsers(ctx context.Context, userID int64, projectName string) ([]*domain.ProjectUser, error) {
	users, err := s.repo.GetProjectUsers(ctx, userID, projectName)
	if err != nil {
		return nil, err
	}

	return users, nil
}
