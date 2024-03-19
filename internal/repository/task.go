package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"rest/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
)

type taskRepo struct {
	db *pgx.Conn
}

func NewTaskRepo(db *pgx.Conn) *taskRepo {
	return &taskRepo{db}
}

func (r *taskRepo) CreateTask(ctx context.Context, authorID int64, task *domain.Task) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	log.Println("task repo create", authorID, task)

	const queryTask = "INSERT INTO tasks (user_id, name, description) VALUES ($1, $2, $3) RETURNING id"

	stmt, err := r.db.Prepare(ctx, "inserTask", queryTask)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var taskID int64
	err = r.db.QueryRow(ctx, stmt.Name, task.Name, task.Description).Scan(&taskID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	const queryUsersTasks = "INSERT INTO users_tasks (user_id, task_id, created_at, expired_at) VALUES ($1, $2, $3, $4)"

	stmt, err = r.db.Prepare(ctx, "inserTaskUser", queryUsersTasks)
	log.Println("[task CREATE repo]", task, err)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	expTime := time.Now().Add(time.Minute) // add mins from exp to curr time

	commandTag, err := r.db.Exec(ctx, stmt.Name, authorID, taskID, time.Now(), expTime)

	// log.Println("[task CREATE repo]", task, expTime, commandTag, err)

	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("create task transaction is not comleted")
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
