package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	_ "github.com/lib/pq"
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

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to open database connection", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()
	if err := db.PingContext(pingCtx); err != nil {
		logger.Error("failed to connect to database", err)
		os.Exit(1)
	}

	router := fiber.New(fiber.Config{
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// TODO: registrar rotas/handlers aqui

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
