package main

import (
	"log"

	// Importaciones internas usando el módulo de go.mod
	"github.com/Manuelgruiz/chamba-ai/apps/api-gateway/db"
	"github.com/Manuelgruiz/chamba-ai/apps/api-gateway/handlers"
	"github.com/Manuelgruiz/chamba-ai/apps/api-gateway/middleware"

	// Paquetes compartidos y librerías externas
	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Inicializar conexión a PostgreSQL local
	db.InitDB()

	app := fiber.New(fiber.Config{
		AppName: "ChambaAI - API Gateway",
	})

	// Middleware global
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Inicializar scraper worker
	worker := provider.NewScraperWorker("default")
	searchHandler := handlers.NewSearchHandler(worker)

	// Grupo de APIs v1
	api := app.Group("/api/v1")

	// --- Rutas Públicas de Autenticación ---
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
	auth.Post("/logout", handlers.Logout)

	// --- Rutas Públicas de Búsqueda ---
	api.Get("/jobs/search", searchHandler.SearchJobs)

	// --- Rutas Protegidas (Requieren Token JWT) ---
	protected := api.Group("/user", middleware.Protected())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		email := c.Locals("email")
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Bienvenido a tu área privada",
			"user_id": userID,
			"email":   email,
		})
	})

	log.Println("API Gateway ejecutándose en http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
