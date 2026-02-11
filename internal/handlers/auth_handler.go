package handlers

import (
	"log"
	"net/http"

	"github.com/echo/internal/models"
	"github.com/echo/internal/service"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService service.AuthService
	//validator   echo.Validator
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		//validator:   validator,
	}
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req models.RegisterRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, "Error binding request")
		}

		resp, err := h.authService.RegisterUser(c.Request().Context(), &req)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, resp)
	}
}
