package repository

import (
	"praktikum/model"

	"github.com/labstack/echo/v4"
)

type Database interface {
	GetAll() ([]model.User, error)
	SaveUser(c echo.Context) error
	Login(c echo.Context) (interface{}, error)
}
