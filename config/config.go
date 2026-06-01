package config

import (
	"os"
	"time"
)

type Config struct {
	DB     DBConfig
	Auth   AuthConfig
	Server ServerConfig
	Logger LoggerConfig
}

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	QueryTimeout time.Duration
	TxTimeout    time.Duration
}

type AuthConfig struct {
	JWTSecret         string
	JWTIssuer         string
	JWTAccessTokenTTL string
	GoogleClientID    string
}

type ServerConfig struct {
	RequestTimeout  time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type LoggerConfig struct {
	Level string
}

func Load() Config {
	return Config{
		DB: DBConfig{
			Host:         os.Getenv("DB_HOST"),
			Port:         os.Getenv("DB_PORT"),
			User:         os.Getenv("DB_USER"),
			Password:     os.Getenv("DB_PASSWORD"),
			Name:         os.Getenv("DB_NAME"),
			SSLMode:      os.Getenv("DB_SSLMODE"),
			QueryTimeout: 2 * time.Second,
			TxTimeout:    5 * time.Second,
		},
		Auth: AuthConfig{
			JWTSecret:         os.Getenv("JWT_SECRET"),
			JWTIssuer:         os.Getenv("JWT_ISSUER"),
			JWTAccessTokenTTL: os.Getenv("JWT_ACCESS_TOKEN_TTL"),
			GoogleClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
		},
		Server: ServerConfig{
			RequestTimeout:  3 * time.Second,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Logger: LoggerConfig{
			Level: os.Getenv("LOG_LEVEL"),
		},
	}
}
