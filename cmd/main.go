package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rest/internal/config"
	"rest/internal/handlers"
	"rest/internal/repository"
	"rest/internal/services"
	"syscall"

	"rest/internal/server"

	"github.com/jackc/pgx/v5"
)

var ctx = context.Background()

func main() {
	cfg, err := config.LoadConfig()
	log.Println(cfg, err)
	if err != nil {
		log.Fatalf("error occured while loading config: %v\n", err)
	}

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

	userRepo := repository.NewUserRepo(conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	taskRepo := repository.NewTaskRepo(conn)

	projectRepo := repository.NewProjectRepo(conn)
	projectService := services.NewProjectService(projectRepo, taskRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	routes := handlers.Routes(userHandler, nil, projectHandler)

	httpServer := new(server.Server)
	go func() {
		if err := httpServer.Run(cfg, routes); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("server is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("server is shutting down")

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}
}
