package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vhgomes/azgher/internal/api/handler"
	"github.com/vhgomes/azgher/internal/postgres/db"
	"github.com/vhgomes/azgher/internal/repository"
	"github.com/vhgomes/azgher/internal/service"
	"github.com/vhgomes/azgher/pkg/config"
	"github.com/vhgomes/azgher/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", err)
		os.Exit(1)
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to parse database config", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		logger.Error("failed to create database pool", err)
		os.Exit(1)
	}
	defer pool.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := pool.Ping(pingCtx); err != nil {
		logger.Error("failed to connect to database", err)
		os.Exit(1)
	}

	router := fiber.New(fiber.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	queries := db.New(pool)
	userRepo := repository.NewUserRepo(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router.Post("/users", userHandler.Create)
	router.Get("/users/:id", userHandler.GetById)
	router.Get("/users/email/:email", userHandler.GetByEmail)
	router.Get("/users/google/:google_id", userHandler.ByGoogleId)

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("starting server", zap.String("port", cfg.Port))
		if err := router.Listen(":" + cfg.Port); err != nil {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("server error", err)
		os.Exit(1)
	case sig := <-shutdown:
		logger.Info("shutting down server", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := router.ShutdownWithContext(ctx); err != nil {
		logger.Error("server forced to shutdown", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
