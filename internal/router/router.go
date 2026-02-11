package router

import (
	"github.com/echo/internal/handlers"
	"github.com/echo/internal/pkg/jwt"
	"github.com/labstack/echo/v5"
)

func Health() echo.HandlerFunc {
	return func(c *echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "ok",
			"service": "echo-api",
		})
	}
}

func Setup(
	e *echo.Echo,
	authHandler *handlers.AuthHandler,
	jwtManager *jwt.Manager,
) {
	// API v1 group
	api := e.Group("/api/v1")
	api.GET("/health", Health())

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register())
}
