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

func main() {
	cfg, err := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.DBName,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	if err = conn.Ping(context.Background()); err != nil {
		log.Fatalf("can't ping db: %s", err)
	}

	// r := chi.NewRouter()
	// r.Use(middleware.RequestID)
	// r.Use(middleware.Logger)
	userRepo := repository.NewUserRepository(conn)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := userHandler.Routes()
	log.Println("server is running")
	http.ListenAndServe(":3000", r)
}
