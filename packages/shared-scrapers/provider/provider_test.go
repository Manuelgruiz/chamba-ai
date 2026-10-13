package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers/provider"
)

func TestAllProviders(t *testing.T) {
	providers := map[string]provider.JobProvider{
		"LinkedIn (Playwright)": provider.NewPlaywrightLinkedInProvider("http://localhost:3001"),
		"InfoJobs":              provider.NewInfoJobsProvider(),
		"Indeed":                provider.NewIndeedProvider(),
	}

	query := "Golang"
	location := "Madrid"

	for name, p := range providers {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			defer cancel()

			t.Logf("[%s] Buscando '%s' en '%s'...", name, query, location)
			offers, err := p.Search(ctx, query, location)

			if err != nil {
				t.Logf("❌ [%s] Error: %v", name, err)
				t.Fail()
				return
			}

			t.Logf("✅ [%s] Encontradas %d ofertas.", name, len(offers))
			if len(offers) > 0 {
				sample := offers[0]
				t.Logf("   Muestra: %s | %s | %s", sample.Title, sample.Company, sample.URL)
			} else {
				t.Logf("⚠️ [%s] No se devolvieron ofertas (posible cambio de HTML o anti-bot).", name)
			}
		})
	}
}
