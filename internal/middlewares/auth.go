package middlewares

import (
	"gotickets/internal/auth"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			var tokenString string

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid Authorization format. Expected Bearer <token>"})
				}
				tokenString = parts[1]
			} else if cookie, err := c.Cookie("access_token"); err == nil {
				tokenString = cookie.Value
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing Authorization header or cookie"})
			}
			claims, err := jwtService.ValidateToken(tokenString, false)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			// Add the claims to context so handlers can access the user info
			c.Set("user", claims)

			return next(c)
		}
	}
}
