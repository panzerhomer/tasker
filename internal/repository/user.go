package repository

import (
	"context"
	"fmt"
	"log"
	"rest/internal/domain"

	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	// db *sqlx.DB
	db *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *userRepo {
	return &userRepo{db}
}

func (ur *userRepo) InsertUser(ctx context.Context, email string, password string) (int64, error) {
	const query = "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id"

	// stmt, err := ur.db.Prepare(query)
	// if err != nil {
	// 	return 0, fmt.Errorf("%w", err)
	// }

	stmt, err := ur.db.Prepare(ctx, "inserUser", query)
	if err != nil {
		log.Println("[err insertuser stmt ]", stmt, err)
		return 0, fmt.Errorf("%w", err)
	}

	// result, err := ur.db.ExecEx(ctx, stmt.Name, nil, email, password)
	// if err != nil {
	// 	return 0, fmt.Errorf("%w", err)
	// }

	var id int64
	err = ur.db.QueryRow(context.Background(), stmt.Name, email, password).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	// result, err := stmt.ExecContext(ctx, email, password)
	// if err != nil {
	// 	return 0, fmt.Errorf("%w", err)
	// }

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, fmt.Errorf("%w", err)
	// }

	return id, nil
}

func (ur *userRepo) GetUser(ctx context.Context, email string) (*domain.User, error) {
	const query = "SELECT id, email, password FROM users WHERE email = $1"

	// stmt, err := ur.db.Prepare(query)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w", err)
	// }

	// row := stmt.QueryRowContext(ctx, email)

	var user domain.User
	// err = row.Scan(&user.ID, &user.Email, &user.Password)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w", err)
	// }
	stmt, err := ur.db.Prepare(ctx, "selectUser", query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	err = ur.db.QueryRow(ctx, stmt.Name, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &user, nil
}
