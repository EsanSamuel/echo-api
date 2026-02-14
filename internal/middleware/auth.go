package middleware

import (
	"net/http"
	"strings"

	"github.com/echo/internal/pkg/jwt"
	"github.com/labstack/echo/v5"
)

func Auth(jwtManager *jwt.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusNotFound, map[string]any{
					"Success": false,
					"Error":   "Auth Header not found!",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"Success": false,
					"Error":   "Invalid header!",
				})
			}
			tokenString := parts[1]

			claims, err := jwtManager.VefifyAccessToken(tokenString)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]any{
					"Success": false,
					"Error":   "Error verfying  access token!",
				})
			}

			c.Set("user_id", claims.ID)
			c.Set("user_email", claims.Email)

			return next(c)
		}
	}
}
