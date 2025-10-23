package main

import (
	"context"
	"net/http"
	"slices"

	"github.com/coreos/go-oidc"
	"github.com/labstack/echo/v4"
)

func authorizeRole(requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			idToken := c.Get("idToken").(*oidc.IDToken)
			var claims map[string]any

			if err := idToken.Claims(&claims); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			realmAccess, ok := claims["realm_access"].(map[string]any)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			roles, ok := realmAccess["roles"].([]any)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			userRoles := make([]string, len(roles))
			if slices.Contains(userRoles, requiredRole) {
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden, "forbidden")
		}
	}
}

func main() {
	e := echo.New()

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "http://localhost:8080/realms/develop")
	if err != nil {
		e.Logger.Fatalf("failed to get provider: %v", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: "api-secure"})
	jwtMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			idTokenStr := authHeader[len("Bearer "):]
			idToken, err := verifier.Verify(ctx, idTokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set("idToken", idToken)
			return next(c)
		}
	}

	e.GET("/private", func(c echo.Context) error {
		return c.String(http.StatusOK, "access allow")
	}, jwtMiddleware)

	e.GET("/admin", func(c echo.Context) error {
		return c.String(http.StatusOK, "access allow admin")
	}, jwtMiddleware, authorizeRole("admin"))

	e.GET("/user", func(c echo.Context) error {
		return c.String(http.StatusOK, "access allow user")
	}, jwtMiddleware, authorizeRole("user"))

	e.Logger.Fatal(e.Start(":8001"))
}
