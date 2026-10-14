package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
)

func TestScraperWorker_FetchAll(t *testing.T) {
	worker := provider.NewScraperWorker("http://localhost:3001")

	ctx := context.Background()
	query := "Golang"
	location := "Madrid"

	t.Logf("🚀 Iniciando búsqueda concurrente de '%s' en '%s'...", query, location)
	start := time.Now()

	offers, errs := worker.FetchAll(ctx, query, location, 35*time.Second)
	duration := time.Since(start)

	t.Logf("⏱️ Búsqueda paralela completada en %v", duration)

	if len(errs) > 0 {
		for platform, err := range errs {
			t.Logf("⚠️ Error parcial en [%s]: %v", platform, err)
		}
	}

	t.Logf("✅ Ofertas totales unificadas y deduplicadas: %d", len(offers))

	if len(offers) == 0 {
		t.Fatalf("Error: No se obtuvo ninguna oferta de ningún proveedor.")
	}

	if len(errs) > 0 {
		for platform, err := range errs {
			t.Logf("⚠️ Error en [%s]: %v", platform, err)
		}
	}

	summary := make(map[string]int)
	for _, o := range offers {
		summary[o.Platform]++
	}

	for platform, count := range summary {
		t.Logf("   - %s: %d ofertas", platform, count)
	}
}
