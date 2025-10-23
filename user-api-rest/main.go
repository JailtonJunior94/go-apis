package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com"},
	{ID: 2, Name: "Jane Doe", Email: "jane@example.com"},
}

func main() {
	e := echo.New()

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
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	for _, user := range users {
		if user.ID == id {
			return c.JSON(http.StatusOK, user)
		}
	}
	return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
}

func createUser(c echo.Context) error {
	var newUser User
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	newUser.ID = len(users) + 1
	users = append(users, newUser)
	return c.JSON(http.StatusCreated, newUser)
}
