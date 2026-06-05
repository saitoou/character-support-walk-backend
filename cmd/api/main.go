package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chocoko/character-support-walk-backend/config"
	"github.com/chocoko/character-support-walk-backend/internal/container"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/observability"
	"github.com/chocoko/character-support-walk-backend/internal/logger"
	appmiddlware "github.com/chocoko/character-support-walk-backend/internal/middleware"
	"github.com/chocoko/character-support-walk-backend/internal/router"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {

	ctx := context.Background()

	shutdown, err := observability.InitTracer(ctx)
	if err != nil {
		slog.Error("failed to init tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			slog.Error("failed to shutdown tracer", "error", err)
		}
	}()

	cfg := config.Load()
	appLogger := logger.New(cfg.Logger.Level)
	slog.SetDefault(appLogger)

	db, err := database.NewDB(cfg)
	if err != nil {
		slog.Error("failed to setup database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ttl, err := time.ParseDuration(cfg.Auth.JWTAccessTokenTTL)
	if err != nil {
		slog.Error("failed to convert access token ttl")
		os.Exit(1)
	}

	jwtService := auth.NewJWTService(cfg.Auth.JWTSecret, ttl, cfg.Auth.JWTIssuer)

	e := echo.New()
	e.Use(echomiddleware.Recover())
	e.Use(echootel.NewMiddleware("character-support-walk-api"))
	e.Use(appmiddlware.CustomCORS())
	e.Use(appmiddlware.RateLimiter())
	e.Use(echomiddleware.RequestID())
	e.Use(appmiddlware.RequestLogger())

	diContainer := container.New(db, jwtService, cfg)
	router.RegisterRoutes(e, diContainer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	sc := echo.StartConfig{
		Address:         ":8080",
		GracefulTimeout: 10 * time.Second,
	}
	slog.Info("server started", "port", 8080)

	if err := sc.Start(ctx, e); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}

	slog.Info("server shutdown completed")
}
