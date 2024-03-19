package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rest/internal/domain"

	"github.com/jackc/pgx/v5"
)

type projectRepo struct {
	db *pgx.Conn
}

func NewProjectRepo(db *pgx.Conn) *projectRepo {
	return &projectRepo{db}
}

func (r *projectRepo) CreateProject(ctx context.Context, authorID int64, project domain.Project) (int64, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	const query = "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING project_id"

	stmt1, err := r.db.Prepare(ctx, "insertProject", query)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	var currentProjectId int64
	err = r.db.QueryRow(ctx, stmt1.Name, project.Name, project.Description).Scan(&currentProjectId)
	if err != nil {
		return 0, err
	}

	const newQuery = "INSERT INTO user_projects (user_id, project_id, user_role) VALUES ($1, $2, $3)"

	stmt2, err := r.db.Prepare(ctx, "insertUserProject", newQuery)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	commandTag, err := r.db.Exec(ctx, stmt2.Name, authorID, currentProjectId, domain.AdminRole)
	log.Println("[!!!!!!!!repository]", authorID, "|", project, "|", err)

	if err != nil {
		return 0, err
	}
	if commandTag.RowsAffected() != 1 {
		return 0, errors.New("no row found to insert")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return currentProjectId, nil
}

func (r *projectRepo) JoinProjectByName(ctx context.Context, projectName string, userID int64, role int8) error {
	project, err := r.GetProjectByName(ctx, projectName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = r.IsMember(ctx, userID, project.Name)
	if err == nil {
		return errors.New("user already in project")
	}

	// log.Println("!!!!!!!!!!!!!!!!", projectName, userID, role, err)
	const newQuery = "INSERT INTO user_projects (user_id, project_id, user_role) VALUES ($1, $2, $3)"
	stmt, err := r.db.Prepare(ctx, "insertUserProjectByName", newQuery)
	log.Println("!!!!!!!!!!!!!!!!", projectName, userID, role, err)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	commandTag, err := r.db.Exec(ctx, stmt.SQL, userID, project.ID, role)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to insert")
	}

	return nil
}

func (r *projectRepo) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
	const query = "SELECT project_id, name, description FROM projects WHERE name = $1"

	var proj domain.Project

	stmt, err := r.db.Prepare(ctx, "selectProjectByName", query)
	log.Println("!!!!!!!!!!!!!!!!!!!!!!", err)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = r.db.QueryRow(ctx, stmt.Name, name).Scan(&proj.ID, &proj.Name, &proj.Description)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &proj, nil
}

func (r *projectRepo) GetProjectById(ctx context.Context, projectID int64) (*domain.Project, error) {
	const query = "SELECT project_id, name, description FROM projects WHERE project_id = $1"

	var proj domain.Project

	stmt, err := r.db.Prepare(ctx, "selectProjectById", query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = r.db.QueryRow(ctx, stmt.Name, projectID).Scan(&proj.ID, &proj.Name, &proj.Description)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &proj, nil
}

func (r *projectRepo) GetProjects(ctx context.Context, userID int64) ([]*domain.Project, error) {
	const query = `
		SELECT 
			p.name, p.description
		FROM 
			projects as p
		INNER JOIN 
			user_projects as up
		ON 
			up.project_id = p.project_id
		WHERE 
			up.user_id = $1`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(&project.Name, &project.Description)
		if err != nil {
			return nil, err
		}
		projects = append(projects, &project)
	}

	log.Println(projects)

	return projects, nil
}

func (r *projectRepo) IsMember(ctx context.Context, userID int64, projectName string) (int8, error) {
	proj, err := r.GetProjectByName(ctx, projectName)
	if err != nil {
		return 0, err
	}

	const query = `SELECT user_role FROM user_projects WHERE user_id = $1 AND project_id = $2`

	var role int8

	err = r.db.QueryRow(ctx, query, userID, proj.ID).Scan(&role)
	if err != nil {
		return 0, fmt.Errorf("%w: ", err)
	}

	return role, nil
}

func (r *projectRepo) UpdateProjectMembers(ctx context.Context, userID int64, targetID int64, projectName string, role int8) error {
	role, err := r.IsMember(ctx, userID, projectName)
	if err != nil {
		return err
	}

	if role != domain.AdminRole {
		return errors.New("missed role or invalid user")
	}

	roleTarget, err := r.IsMember(ctx, targetID, projectName)
	if err != nil {
		return errors.New("target user not member")
	}

	if roleTarget == domain.AdminRole {
		return errors.New("target user already admin")
	}

	const query = `UPDATE user_projects SET user_role = $1 WHERE user_id = $2`

	tag, err := r.db.Exec(ctx, query, role, targetID)
	if err != nil {
		return errors.New("can't update user's role")
	}
	if tag.RowsAffected() != 1 {
		return errors.New("can't update user's role")
	}

	return nil
}
