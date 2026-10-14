package main

import (
	"log"

	"github.com/Manuelgruiz/chamba-ai/apps/api-gateway/handlers"
	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "ChambaAI API Gateway",
		ServerHeader: "Fiber",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	scraperWorker := provider.NewScraperWorker("http://localhost:3001")
	searchHandler := handlers.NewSearchHandler(scraperWorker)

	v1 := app.Group("/api/v1")

	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "api-gateway"})
	})

	v1.Get("/jobs/search", searchHandler.SearchJobs)

	log.Println("🚀 ChambaAI API Gateway escuchando en http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error iniciando el API Gateway: %v", err)
	}
}
