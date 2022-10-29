package repository

import (
	mid "praktikum/middleware"
	"praktikum/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GormSql struct {
	DB *gorm.DB
}

func (m *GormSql) GetAll() ([]model.User, error) {
	var users []model.User

	if err := m.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (m *GormSql) SaveUser(c echo.Context) error {
	var user model.User

	c.Bind(&user)

	if err := m.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
func (m *GormSql) Login(c echo.Context) (interface{}, error) {
	var u model.User
	c.Bind(&u)

	err := m.DB.Where("email = ? AND password = ?", u.Email, u.Password).First(&u).Error
	if err != nil {
		return nil, err
	}

	// TODO Create JWT Token
	token, err := mid.CreateToken(int(u.ID), u.Email)
	if err != nil {
		return nil, err
	}
	// TODO Bind created token into user token
	result := map[string]interface{}{
		"user":  u,
		"token": token,
	}

	return result, nil
}

func NewGorm(DB *gorm.DB) Database {
	return &GormSql{
		DB: DB,
	}
}
