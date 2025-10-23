package main

import (
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserDTO struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserFilterDTO struct {
	Name string `query:"name"`
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.BasicAuth(func(user, pass string, ctx echo.Context) (bool, error) {
		if user == "admin" && pass == "password" {
			return true, nil
		}
		return false, nil
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "Echo/4.0")
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.POST("/users", func(c echo.Context) error {
		var user UserDTO
		if err := c.Bind(&user); err != nil {
			return c.JSON(400, map[string]any{"error": "Invalid input"})
		}
		return c.JSON(201, map[string]any{
			"id":    1,
			"name":  user.Name,
			"email": user.Email,
		})
	})

	e.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		return c.JSON(200, map[string]any{
			"id":   id,
			"name": "Jailton Junior",
		})
	})

	e.GET("/users/filter", func(c echo.Context) error {
		var filter UserFilterDTO
		if err := (&echo.DefaultBinder{}).BindQueryParams(c, &filter); err != nil {
			return c.JSON(400, map[string]any{"error": "Invalid query parameters"})
		}
		return c.JSON(200, map[string]any{
			"filter_name": filter.Name,
		})
	})

	e.PUT("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		name := c.FormValue("name")

		avatar, err := c.FormFile("avatar")
		if err != nil {
			return c.JSON(400, map[string]any{"error": "Avatar upload failed"})
		}

		file, err := avatar.Open()
		if err != nil {
			return c.JSON(500, map[string]any{"error": "Failed to open avatar file"})
		}

		dest, err := os.Create("./" + avatar.Filename)
		if err != nil {
			return c.JSON(500, map[string]any{"error": "Failed to create destination file"})
		}
		defer dest.Close()

		_, err = io.Copy(dest, file)
		if err != nil {
			return c.JSON(500, map[string]any{"error": "Failed to save avatar file"})
		}

		return c.JSON(200, map[string]any{
			"id":     id,
			"name":   name,
			"avatar": avatar.Filename,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
