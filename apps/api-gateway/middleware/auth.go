package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("chamba_secret_key_super_segura_local")

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tokenString string

		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = c.Cookies("jwt")
		}

		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{"error": "Acceso no autorizado. Token faltante."})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": "Token inválido o expirado."})
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"])
		c.Locals("email", claims["email"])

		return c.Next()
	}
}
