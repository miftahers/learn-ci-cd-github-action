package controller

import (
	"net/http"
	"praktikum/services"

	"github.com/labstack/echo/v4"
)

type HandlerUser struct {
	UserServices services.UserServices
}

func (h *HandlerUser) LoginUser(c echo.Context) error {
	result, err := h.UserServices.LoginUser(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Fail login user",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success",
		"data":    result,
	})
}

func (h *HandlerUser) GetAllUsers(c echo.Context) error {
	users, err := h.UserServices.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Success",
		"data":    users,
	})

}

func (h *HandlerUser) CreateUser(c echo.Context) error {
	err := h.UserServices.CreateUser(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Created!",
	})
}
