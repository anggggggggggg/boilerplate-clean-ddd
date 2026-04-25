package middleware

import (
	"strings"
	"time"

	"github.com/anggggggggggg/boilerplate-clean-and-ddd/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("test") // Ganti dengan secret dari env

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "Missing or invalid Authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return response.WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "Invalid token claims")
		}

		userID, ok := claims["user_id"]
		if !ok {
			return response.WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "user_id not found in token")
		}
		c.Locals("user_id", userID)

		return c.Next()
	}
}

func GenerateToken(userID string, username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
