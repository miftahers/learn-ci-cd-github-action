package routes

import (
	"praktikum/config"
	"praktikum/controller"
	"praktikum/repository"
	userv "praktikum/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func Init(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// TODO handle TrailingSlash
	e.Pre(middleware.RemoveTrailingSlash())

	// TODO Implement logger
	middleware.Logger()

	repoUser := repository.NewGorm(db)
	services := userv.NewUser(repoUser)
	handlerUser := controller.HandlerUser{
		UserServices: services,
	}
	e.GET("/users", handlerUser.GetAllUsers, middleware.JWT([]byte(config.TokenSecret)))
	e.POST("/users/login", handlerUser.LoginUser)
	e.POST("/users", handlerUser.CreateUser)

	return e
}
