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

func (h *AuthHandler) VerifyUser() echo.HandlerFunc {
	return func(c *echo.Context) error {

		token := c.QueryParam("token")

		err := h.authService.VerifyEmail(c.Request().Context(), token)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusCreated, "User verified!")
	}
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c *echo.Context) error {
		var req models.LoginRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, "Error binding request")
		}

		user, err := h.authService.Login(c.Request().Context(), &req)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		accessCookie := &http.Cookie{
			Name:     "access_token",
			Value:    user.AccessToken,
			Path:     "/",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}

		refreshCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    user.RefreshToken,
			Path:     "/",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		}

		c.SetCookie(accessCookie)
		c.SetCookie(refreshCookie)

		return c.JSON(http.StatusCreated, user)
	}
}
