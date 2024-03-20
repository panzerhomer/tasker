package repository

import (
	"context"
	"errors"
	"fmt"
	"rest/internal/domain"

	"github.com/jackc/pgx/v5"
)

type taskRepo struct {
	db *pgx.Conn
}

func NewTaskRepo(db *pgx.Conn) *taskRepo {
	return &taskRepo{db}
}

func (r *taskRepo) AssignUserToTask(ctx context.Context, userID int64, projectID int64, task domain.Task) error {
	const query = `INSERT INTO tasks 
		(name, description, status, deadline, project_id, assigned_user_id, author_id) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	tag, err := r.db.Exec(ctx, query, task.Name, task.Description, task.Deadline, projectID, task.AssignedUserID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() != 1 {
		return errors.New("can't insert task into project")
	}

	return nil
}

func (r *taskRepo) UpdateTask(ctx context.Context, userID int64, status int64) error {
	const query = "UPDATE tasks SET status = $1 WHERE assigned_user_id = userID"

	stmt, err := r.db.Prepare(ctx, "updateTask", query)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = r.db.Exec(ctx, stmt.Name, status, userID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (r *taskRepo) DeleteTask(ctx context.Context, authorID int64, taskID int64) error {
	const query = "DELETE FROM tasks WHERE id = $1 and user_id = $2"

	stmt, err := r.db.Prepare(ctx, "deleteTask", query)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	_, err = r.db.Exec(ctx, stmt.Name, taskID, authorID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
