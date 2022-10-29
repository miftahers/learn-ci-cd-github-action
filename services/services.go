package services

import (
	"praktikum/dto"
	"praktikum/repository"

	"github.com/labstack/echo/v4"
)

type UserServices interface {
	GetAllUsers() ([]dto.User, error)
	LoginUser(c echo.Context) (interface{}, error)
	CreateUser(c echo.Context) error
}

type userServices struct {
	repository.Database
}

func (us *userServices) GetAllUsers() ([]dto.User, error) {
	res, err := us.GetAll()
	if err != nil {
		return nil, err
	}

	var dtos []dto.User

	for _, v := range res {
		dtos = append(dtos, dto.User{
			Email: v.Email,
		})
	}

	return dtos, nil
}

func (us *userServices) LoginUser(c echo.Context) (interface{}, error) {

	result, err := us.Login(c)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (us *userServices) CreateUser(c echo.Context) error {
	err := us.SaveUser(c)
	if err != nil {
		return err
	}

	return nil
}

func NewUser(repo repository.Database) UserServices {
	return &userServices{
		repo,
	}
}
