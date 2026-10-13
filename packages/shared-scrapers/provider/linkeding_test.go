package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
)

func TestPlaywrightLinkedInProvider_Search(t *testing.T) {
	// Instanciar apuntando al servicio local de Node/TypeScript
	p := provider.NewPlaywrightLinkedInProvider("http://localhost:3001")

	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	query := "Golang"
	location := "Spain"

	t.Logf("Buscando ofertas de '%s' en '%s' vía Playwright...", query, location)
	offers, err := p.Search(ctx, query, location)

	if err != nil {
		t.Fatalf("Error durante la búsqueda en Playwright: %v", err)
	}

	if len(offers) == 0 {
		t.Log("Aviso: No se encontraron ofertas en la búsqueda.")
		return
	}

	first := offers[0]
	t.Logf("¡Éxito! Se obtuvieron %d ofertas de LinkedIn.", len(offers))
	t.Logf("Muestra de la primera oferta:")
	t.Logf("- Título: %s", first.Title)
	t.Logf("- Empresa: %s", first.Company)
	t.Logf("- Ubicación: %s", first.Location)
	t.Logf("- URL: %s", first.URL)

	if first.Title == "" || first.URL == "" {
		t.Errorf("Inconsistencia en los datos: Título o URL vacíos")
	}
}
