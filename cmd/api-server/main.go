package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Nicholas2012/time-tracker/docs"
	"github.com/Nicholas2012/time-tracker/internal/api"
	"github.com/Nicholas2012/time-tracker/internal/config"
	"github.com/Nicholas2012/time-tracker/internal/repository"
	"github.com/Nicholas2012/time-tracker/internal/usecase"
	"github.com/Nicholas2012/time-tracker/pkg/database"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	slog.Info("Starting server...")

	config := config.New()

	db, err := database.New(config.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.ApplyMigrations(db); err != nil {
		slog.Error("Failed to apply migrations", "error", err)
		os.Exit(1)
	}

	repo := repository.New(db)
	svc := usecase.New(repo)
	api := api.New(svc)

	api.AddRoutes(http.DefaultServeMux)
	http.Handle("/swagger/", httpSwagger.Handler())

	slog.Info("Server started", "listen", config.Listen)
	if err := http.ListenAndServe(config.Listen, nil); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	// todo add graceful shutdown
}
