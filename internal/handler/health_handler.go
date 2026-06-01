package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v5"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *echo.Context) error {

	var result int

	err := h.db.QueryRow("SELECT 1").Scan(&result)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "db error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
		"db":     result,
	})
}
