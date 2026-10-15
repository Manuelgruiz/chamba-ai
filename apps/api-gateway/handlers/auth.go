package handlers

import (
	"context"
	"time"

	"github.com/Manuelgruiz/chamba-ai/apps/api-gateway/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("chamba_secret_key_super_segura_local")

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cuerpo de petición inválido"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al procesar la contraseña"})
	}

	query := `INSERT INTO users (email, password_hash, full_name) VALUES ($1, $2, $3) RETURNING id`
	var userID string
	err = db.DB.QueryRow(context.Background(), query, req.Email, string(hashedPassword), req.FullName).Scan(&userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "El usuario ya existe o hubo un problema al crearlo"})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"user_id": userID,
		"message": "Usuario registrado exitosamente",
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Petición inválida"})
	}

	var storedHash, userID, fullName string
	query := `SELECT id, password_hash, full_name FROM users WHERE email = $1`
	err := db.DB.QueryRow(context.Background(), query, req.Email).Scan(&userID, &storedHash, &fullName)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales incorrectas"})
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   req.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "No se pudo generar el token JWT"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{
		"success": true,
		"token":   tokenString,
		"user": fiber.Map{
			"id":        userID,
			"email":     req.Email,
			"full_name": fullName,
		},
	})
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})
	return c.JSON(fiber.Map{"success": true, "message": "Sesión cerrada correctamente"})
}
