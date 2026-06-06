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
	echoprometheus "github.com/labstack/echo-contrib/v5/echoprometheus"
	echootel "github.com/labstack/echo-opentelemetry"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		slog.Error("application error", "error", err)
		os.Exit(1)
	}
}

func run() error {

	ctx := context.Background()

	shutdown, err := observability.InitTracer(ctx)
	if err != nil {
		return err
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
		return err
	}
	defer db.Close()

	ttl, err := time.ParseDuration(cfg.Auth.JWTAccessTokenTTL)
	if err != nil {
		return err
	}

	jwtService := auth.NewJWTService(cfg.Auth.JWTSecret, ttl, cfg.Auth.JWTIssuer)

	e := echo.New()
	e.Use(echomiddleware.Recover())
	e.Use(echootel.NewMiddleware("character-support-walk-api"))
	e.Use(echoprometheus.NewMiddleware("character_support_walk"))
	e.Use(appmiddlware.CustomCORS())
	e.Use(appmiddlware.RateLimiter())
	e.Use(echomiddleware.RequestID())
	e.Use(appmiddlware.RequestLogger())

	e.GET("/metrics", echoprometheus.NewHandler())

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
		return err
	}

	slog.Info("server shutdown completed")
	return nil
}
