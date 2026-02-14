package handlers

import (
	"log"
	"net/http"

	"github.com/echo/internal/service"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userService service.UserService
	//validator   echo.Validator
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		//validator:   validator,
	}
}

func (h *UserHandler) FindUser() echo.HandlerFunc {
	return func(c *echo.Context) error {
		id, ok := c.Get("user_id").(int32)
		if !ok {
			return c.JSON(http.StatusInternalServerError, "UserID not found")
		}

		resp, err := h.userService.FindUser(c.Request().Context(), id)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}
