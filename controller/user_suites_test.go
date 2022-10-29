package controller

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"praktikum/dto"
	"praktikum/model"
	"praktikum/repository"
	"praktikum/services"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type suiteUsers struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	handler HandlerUser
}
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestSuiteUser(t *testing.T) {
	suite.Run(t, new(suiteUsers))
}

func (s *suiteUsers) SetupSuite() {
	// Create mock db
	db, mock, err := sqlmock.New()
	s.NoError(err)

	var gDB *gorm.DB
	gDB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	s.NoError(err)

	repoUser := repository.NewGorm(gDB)
	services := services.NewUser(repoUser)
	handlerUser := HandlerUser{
		UserServices: services,
	}

	s.handler = handlerUser
	s.mock = mock
}

func (s *suiteUsers) TearDownSuite() {
	s.mock = nil
}

func (s *suiteUsers) TestCreateUser() {
	testCase := []struct {
		name          string
		path          string
		method        string
		expectCode    int
		Body          model.User
		expectMessage string
	}{
		{
			name:       "normal create user",
			path:       "/users",
			method:     http.MethodPost,
			expectCode: http.StatusCreated,
			Body: model.User{
				Model: &gorm.Model{
					ID: 1,
				},
				Email:    "superb@test.com",
				Password: "superb123",
			},
			expectMessage: "Created!",
		},
	}

	for _, v := range testCase {
		s.T().Run(v.name, func(t *testing.T) {
			s.mock.ExpectBegin()
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`email`,`password`) VALUES (?,?,?,?,?)")).
				WithArgs(AnyTime{}, AnyTime{}, nil, "", "").
				WillReturnResult(sqlmock.NewResult(1, 1))
			s.mock.ExpectCommit()

			res, _ := json.Marshal(v.Body)
			request := httptest.NewRequest(v.method, "/", bytes.NewBuffer(res))
			writter := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(request, writter)
			ctx.SetPath(v.path)

			if s.NoError(s.handler.CreateUser(ctx)) {
				body := writter.Body.Bytes()

				type Response struct {
					Message string `json:"message"`
				}
				var resp Response
				if err := json.Unmarshal(body, &resp); err != nil {
					s.Error(err, "error unmarshalling")
				}

				s.Equal(v.expectCode, writter.Result().StatusCode)
				s.Equal(v.expectMessage, resp.Message)
			}
		})
	}
}

func (s *suiteUsers) TestCreateUserError() {
	testCase := []struct {
		name          string
		path          string
		method        string
		expectCode    int
		Body          model.User
		expectMessage string
	}{
		{
			name:       "error create user",
			path:       "/users",
			method:     http.MethodPost,
			expectCode: http.StatusInternalServerError,
			Body: model.User{
				Model: &gorm.Model{
					ID: 1,
				},
				Email:    "superb@test.com",
				Password: "superb123",
			},
			expectMessage: "Internal Server Error",
		},
	}

	for _, v := range testCase {
		s.T().Run(v.name, func(t *testing.T) {
			s.mock.ExpectBegin()
			s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`email`,`password`) VALUES (?,?,?,?,?)")).
				WithArgs(AnyTime{}, AnyTime{}, nil, "", "").
				WillReturnError(errors.New("Internal Server Error"))
			s.mock.ExpectRollback()

			res, _ := json.Marshal(v.Body)
			request := httptest.NewRequest(v.method, "/", bytes.NewBuffer(res))
			writter := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(request, writter)
			ctx.SetPath(v.path)

			if s.NoError(s.handler.CreateUser(ctx)) {
				body := writter.Body.Bytes()

				type Response struct {
					Message string `json:"message"`
				}
				var resp Response
				if err := json.Unmarshal(body, &resp); err != nil {
					s.Error(err, "error unmarshalling")
				}

				s.Equal(v.expectCode, writter.Result().StatusCode)
				s.Equal(v.expectMessage, resp.Message)
			}
		})
	}
}

func (s *suiteUsers) TestGetAllUsers() {
	testCase := []struct {
		name          string
		path          string
		method        string
		expectCode    int
		expectBody    []dto.User
		expectMessage string
	}{
		{
			name:       "normal get users",
			path:       "/users",
			method:     http.MethodGet,
			expectCode: http.StatusOK,
			expectBody: []dto.User{
				{
					Email: "superb@test.com",
				},
				{
					Email: "amazing@test.com",
				},
			},
			expectMessage: "Success",
		},
	}

	for _, v := range testCase {
		s.T().Run(v.name, func(t *testing.T) {
			expectedRows := s.mock.NewRows([]string{"email"}).
				AddRow("superb@test.com").
				AddRow("amazing@test.com")
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
				WillReturnRows(expectedRows)

			r := httptest.NewRequest(v.method, "/", nil)
			w := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath(v.path)

			if s.NoError(s.handler.GetAllUsers(ctx)) {
				body := w.Body.Bytes()

				type response struct {
					Message string     `json:"message"`
					Data    []dto.User `json:"data"`
				}
				var resp response
				if err := json.Unmarshal(body, &resp); err != nil {
					s.Error(err, "error unmarshalling")
				}

				s.Equal(v.expectCode, w.Result().StatusCode)
				s.Equal(v.expectBody, resp.Data)
			}
		})
	}
}

func (s *suiteUsers) TestGetAllUsersError() {
	testCase := []struct {
		name         string
		path         string
		method       string
		expectResult string
	}{
		{
			name:         "error get users",
			path:         "/users",
			method:       http.MethodGet,
			expectResult: echo.NewHTTPError(http.StatusInternalServerError, "Record Not found").Error(),
		},
	}
	for _, v := range testCase {
		s.T().Run(v.name, func(t *testing.T) {
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL")).
				WillReturnError(errors.New("Record Not Found"))
			r := httptest.NewRequest(v.method, "/", nil)
			w := httptest.NewRecorder()
			e := echo.New()
			ctx := e.NewContext(r, w)
			ctx.SetPath(v.path)
			err := s.handler.GetAllUsers(ctx)
			s.Equal(v.expectResult, err.Error())
		})
	}
}

func (s *suiteUsers) TestLoginUser() {
	testCase := []struct {
		name          string
		path          string
		method        string
		body          model.User
		expectCode    int
		expectMessage string
	}{
		{
			name:   "normal login user",
			path:   "/users/login",
			method: http.MethodPost,
			body: model.User{
				Email:    "orang@test.com",
				Password: "orang123",
			},
			expectCode:    http.StatusOK,
			expectMessage: "Success",
		},
	}

	for _, v := range testCase {
		s.T().Run(v.name, func(t *testing.T) {
			expectRows := s.mock.NewRows([]string{"id", "email", "password"}).AddRow(1, "orang@test.com", "orang123")
			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ? AND password = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
				WithArgs("", "").
				WillReturnRows(expectRows)

			res, _ := json.Marshal(v.body)
			r := httptest.NewRequest(v.method, "/", bytes.NewBuffer(res))
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)
			ctx.SetPath(v.path)

			if s.NoError(s.handler.LoginUser(ctx)) {
				body := w.Body.Bytes()

				type response struct {
					Message string      `json:"message"`
					Data    interface{} `json:"data"`
				}
				var resp response
				err := json.Unmarshal(body, &resp)
				if err != nil {
					s.Error(err)
				}

				s.Equal(v.expectCode, w.Result().StatusCode)
				s.Equal(v.expectMessage, resp.Message)
			}
		})
	}
}

// func (s *suiteUsers) TestLoginUserError() {
// 	testCase := []struct {
// 		name        string
// 		path        string
// 		method      string
// 		body        model.User
// 		expectError string
// 	}{
// 		{
// 			name:   "error login user",
// 			path:   "/users/login",
// 			method: http.MethodPost,
// 			body: model.User{
// 				Email:    "orang@test.com",
// 				Password: "orang123",
// 			},
// 			expectError: echo.NewHTTPError(http.StatusInternalServerError, "Record Not found").Error(),
// 		},
// 	}
// 	for _, v := range testCase {
// 		s.T().Run(v.name, func(t *testing.T) {
// 			s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE (email = ? AND password = ?) AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")).
// 				WillReturnError(errors.New("Fail login user"))
// 			res, _ := json.Marshal(v.body)
// 			r := httptest.NewRequest(v.method, "/", bytes.NewBuffer(res))
// 			w := httptest.NewRecorder()
// 			ctx := echo.New().NewContext(r, w)
// 			ctx.SetPath(v.path)
// 			err := (s.handler.LoginUser(ctx)).Error()
// 			s.Equal(v.expectError, err)
// 		})
// 	}
// }
