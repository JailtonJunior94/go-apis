package main

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func customerHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	msg := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message.(string)
	}
	c.JSON(code, map[string]string{"error": msg})
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.HTTPErrorHandler = customerHTTPErrorHandler

	e.GET("/users", getUsers)
	e.GET("/users/:id", getUser)
	e.POST("/users", createUser)

	e.Logger.Fatal(e.Start(":8080"))
}

func getUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, users)
}

func getUser(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	for _, user := range users {
		if user.ID == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return echo.NewHTTPError(http.StatusNotFound, "User not found")
}

func createUser(c echo.Context) error {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	if err := c.Validate(&newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser.ID = len(users) + 1
	users = append(users, newUser)
	return c.JSON(http.StatusCreated, newUser)
}
