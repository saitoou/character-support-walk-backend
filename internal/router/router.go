package router

import (
	"github.com/chocoko/character-support-walk-backend/internal/container"
	"github.com/labstack/echo/v5"
)

func RegisterRoutes(
	e *echo.Echo,
	c *container.Container,
) {

	jwtAuth := c.JWTAuth

	// health check
	e.GET("/health", c.HealthHandler.Health)

	apiV1 := e.Group("/api/v1")

	// master
	apiV1.GET("/walk-options", c.WalkOptionHandler.GetWalkOptions, jwtAuth)

	// aggregation
	apiV1.GET("/home", c.HomeHandler.GetHome, jwtAuth)

	// walk
	apiV1.POST("/walks", c.WalkHandler.StartWalk, jwtAuth)
	apiV1.GET("/walks", c.WalkHandler.GetWalks, jwtAuth)
	apiV1.GET("/walks/:walkID", c.WalkHandler.GetWalk, jwtAuth)
	apiV1.PATCH("/walks/:walkID/complete", c.WalkHandler.CompleteWalk, jwtAuth)
	apiV1.PATCH("/walks/:walkID/cancel", c.WalkHandler.CancelWalk, jwtAuth)

	// me
	apiV1.GET("/walker", c.UserHandler.GetUser, jwtAuth)
	apiV1.PATCH("/walker", c.UserHandler.UpdateWalker, jwtAuth)
	apiV1.DELETE("/walker", c.UserHandler.DeactivateWalker, jwtAuth)

	// auth
	apiV1.POST("/auth/dev-login", c.AuthHandler.DevLogin)
	apiV1.POST("/auth/google", c.AuthHandler.GoogleLogin)
	apiV1.POST("/auth/refresh", c.AuthHandler.Refresh)
	apiV1.POST("/auth/logout", c.AuthHandler.Logout)

}
