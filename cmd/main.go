package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest/internal/config"
	"rest/internal/handlers"
	"rest/internal/repository"
	"rest/internal/services"

	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

func main() {
	cfg, _ := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("can't ping db: %s", err)
	}

	//conn.Exec(ctx, "DELETE FROM users")

	userRepo := repository.NewUserRepo(conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	taskRepo := repository.NewTaskRepo(conn)
	taskService := services.NewTaskService(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	projectRepo := repository.NewProjectRepo(conn)
	projectService := services.NewProjectService(projectRepo, taskRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	r := handlers.Routes(userHandler, taskHandler, projectHandler)
	log.Println("server is running")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal("server stopped with ", err)
	}
}
