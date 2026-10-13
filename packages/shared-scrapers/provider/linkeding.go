package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PlaywrightLinkedInProvider struct {
	ServiceURL string
	Client     *http.Client
}

func NewPlaywrightLinkedInProvider(serviceURL string) *PlaywrightLinkedInProvider {
	if serviceURL == "" {
		serviceURL = "http://localhost:3001"
	}
	return &PlaywrightLinkedInProvider{
		ServiceURL: serviceURL,
		Client: &http.Client{
			Timeout: 40 * time.Second,
		},
	}
}

type scrapeResponse struct {
	Success bool       `json:"success"`
	Count   int        `json:"count"`
	Offers  []JobOffer `json:"offers"`
	Error   string     `json:"error,omitempty"`
}

func (p *PlaywrightLinkedInProvider) Search(ctx context.Context, query, location string) ([]JobOffer, error) {
	payload := map[string]string{
		"query":    query,
		"location": location,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", p.ServiceURL+"/scrape", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error preparando petición a microservicio: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("microservicio de LinkedIn no disponible en %s: %w", p.ServiceURL, err)
	}
	defer resp.Body.Close()

	var result scrapeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Playwright: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("error en scraping Playwright: %s", result.Error)
	}

	return result.Offers, nil
}
