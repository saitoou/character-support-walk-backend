package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/chocoko/character-support-walk-backend/config"
)

func NewDB(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	slog.Info("connecting database...")
	slog.Debug("connecting database",
		"host", cfg.DB.Host,
		"port", cfg.DB.Port,
		"user", cfg.DB.User,
		"dbname", cfg.DB.Name,
		"sslmode", cfg.DB.SSLMode,
	)

	database, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql open failed: %w", err)
	}

	if err := database.Ping(); err != nil {
		return nil, fmt.Errorf("datanase ping failed: %w", err)
	}

	slog.Info("database connected")

	return database, nil

}
