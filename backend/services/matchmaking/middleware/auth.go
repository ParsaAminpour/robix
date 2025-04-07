package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthConfig struct {
	JWTSecret      string
	UserServiceURL string
}

func AuthMiddleware(config AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			// Extract the token from the Authorization header
			// Format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(config.JWTSecret), nil
			})

			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// Add user info to context
				c.Set("user_id", claims["user_id"])
				c.Set("username", claims["username"])
				return next(c)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
		}
	}
}
