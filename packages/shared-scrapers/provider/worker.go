package provider

import (
	"context"
	"sync"
	"time"
)

// ScrapingResult contiene el resultado de un provider individual.
type ScrapingResult struct {
	Platform string
	Offers   []JobOffer
	Err      error
}

// ScraperWorker es el orquestador que gestiona la ejecución paralela de los proveedores.
type ScraperWorker struct {
	Providers map[string]JobProvider
}

// NewScraperWorker inicializa el orquestador con todos los scrapers registrados.
func NewScraperWorker(playwrightURL string) *ScraperWorker {
	if playwrightURL == "" {
		playwrightURL = "http://localhost:3001"
	}

	return &ScraperWorker{
		Providers: map[string]JobProvider{
			"linkedin": NewPlaywrightLinkedInProvider(playwrightURL),
			"infojobs": NewInfoJobsProvider(),
			"indeed":   NewIndeedProvider(),
		},
	}
}

// FetchAllExecutes la búsqueda en todos los scrapers en paralelo.
// Aplica un timeout global y consolida las ofertas eliminando duplicados.
func (w *ScraperWorker) FetchAll(ctx context.Context, query, location string, globalTimeout time.Duration) ([]JobOffer, map[string]error) {
	ctx, cancel := context.WithTimeout(ctx, globalTimeout)
	defer cancel()

	numProviders := len(w.Providers)
	resultsChan := make(chan ScrapingResult, numProviders)
	var wg sync.WaitGroup

	// Lanza una Goroutine por cada scraper
	for name, p := range w.Providers {
		wg.Add(1)
		go func(platform string, provider JobProvider) {
			defer wg.Done()

			offers, err := provider.Search(ctx, query, location)
			resultsChan <- ScrapingResult{
				Platform: platform,
				Offers:   offers,
				Err:      err,
			}
		}(name, p)
	}

	// Cerrar el canal cuando todas las Goroutines hayan finalizado
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var allOffers []JobOffer
	errorsMap := make(map[string]error)
	seenURLs := make(map[string]bool)

	// Procesar los resultados a medida que llegan
	for res := range resultsChan {
		if res.Err != nil {
			errorsMap[res.Platform] = res.Err
			continue
		}

		// Consolidar y deduplicar ofertas por URL
		for _, offer := range res.Offers {
			if offer.URL != "" && !seenURLs[offer.URL] {
				seenURLs[offer.URL] = true
				allOffers = append(allOffers, offer)
			}
		}
	}

	return allOffers, errorsMap
}
