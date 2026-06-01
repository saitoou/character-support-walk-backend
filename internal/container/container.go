package container

import (
	"database/sql"

	"github.com/chocoko/character-support-walk-backend/config"
	"github.com/chocoko/character-support-walk-backend/internal/handler"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/auth"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/database"
	"github.com/chocoko/character-support-walk-backend/internal/infrastructure/repository"
	"github.com/chocoko/character-support-walk-backend/internal/middleware"
	"github.com/chocoko/character-support-walk-backend/internal/usecase"
	"github.com/labstack/echo/v5"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Container struct {
	JWTAuth echo.MiddlewareFunc

	HealthHandler     *handler.HealthHandler
	WalkOptionHandler *handler.WalkOptionHandler
	WalkHandler       *handler.WalkHandler
	UserHandler       *handler.UserHandler
	HomeHandler       *handler.HomeHandler
	AuthHandler       *handler.AuthHandler
}

func New(
	sqlDB *sql.DB,
	jwtService *auth.JWTService,
	cfg config.Config,
) *Container {
	jwtAuth := middleware.JWTAuthMiddleware(jwtService)
	tx := database.NewTransaction(sqlDB)
	idTokenValidator := auth.NewGoogleIDTokenVerifier(
		cfg.Auth.GoogleClientID,
	)

	// Repository
	walkOptionRepo := repository.NewWalkOptionRepository(sqlDB)
	walkRepo := repository.NewWalkRepository(sqlDB)
	userRepo := repository.NewUserRepository(sqlDB)
	authRepo := repository.NewAuthRepository(sqlDB)
	characterRepo := repository.NewCharacterRepository(sqlDB)

	// Usecase
	walkOptionUsecase := usecase.NewWalkOptionUsecase(walkOptionRepo)
	walkUsecase := usecase.NewWalkUsecase(walkRepo, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, characterRepo, authRepo, tx)
	homeUsecase := usecase.NewHomeUsecase(userRepo, characterRepo, walkRepo)
	authUsecase := usecase.NewAuthUsecase(
		authRepo,
		userRepo,
		characterRepo,
		jwtService,
		idTokenValidator,
		tx,
	)

	// Handler
	healthHandler := handler.NewHealthHandler(sqlDB)
	walkOptionHandler := handler.NewWalkOptionHandler(walkOptionUsecase)
	walkHandler := handler.NewWalkHandler(walkUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	homeHandler := handler.NewHomeHandler(homeUsecase, cfg.Server.RequestTimeout)
	authHandler := handler.NewAuthHandler(authUsecase, jwtService)

	return &Container{
		JWTAuth:           jwtAuth,
		HealthHandler:     healthHandler,
		WalkOptionHandler: walkOptionHandler,
		WalkHandler:       walkHandler,
		UserHandler:       userHandler,
		HomeHandler:       homeHandler,
		AuthHandler:       authHandler,
	}

}
