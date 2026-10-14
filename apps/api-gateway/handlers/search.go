package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
	"github.com/gofiber/fiber/v2"
)

type SearchHandler struct {
	Worker *provider.ScraperWorker
}

func NewSearchHandler(worker *provider.ScraperWorker) *SearchHandler {
	return &SearchHandler{
		Worker: worker,
	}
}

// SearchJobs maneja la petición GET /api/v1/jobs/search
func (h *SearchHandler) SearchJobs(c *fiber.Ctx) error {
	query := c.Query("q")
	location := c.Query("l")

	if query == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "El parámetro de búsqueda 'q' es obligatorio.",
		})
	}

	// Contexto con timeout global para la petición HTTP del usuario
	ctx, cancel := context.WithTimeout(c.UserContext(), 40*time.Second)
	defer cancel()

	start := time.Now()
	offers, errs := h.Worker.FetchAll(ctx, query, location, 35*time.Second)
	duration := time.Since(start)

	// Mapear errores parciales si algún proveedor falló
	partialErrors := make(map[string]string)
	for platform, err := range errs {
		partialErrors[platform] = err.Error()
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"query":         query,
		"location":      location,
		"total_results": len(offers),
		"latency_ms":    duration.Milliseconds(),
		"errors":        partialErrors,
		"data":          offers,
	})
}
