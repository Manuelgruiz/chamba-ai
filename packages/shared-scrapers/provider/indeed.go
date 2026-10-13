package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IndeedProvider struct {
	ServiceURL string
	Client     *http.Client
}

func NewIndeedProvider() *IndeedProvider {
	return &IndeedProvider{
		ServiceURL: "http://localhost:3001",
		Client: &http.Client{
			Timeout: 40 * time.Second,
		},
	}
}

func (p *IndeedProvider) Search(ctx context.Context, query, location string) ([]JobOffer, error) {
	payload := map[string]string{
		"query":    query,
		"location": location,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, "POST", p.ServiceURL+"/scrape-indeed", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creando petición para Indeed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("el servicio Playwright no respondió para Indeed: %w", err)
	}
	defer resp.Body.Close()

	var result scrapeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta de Indeed: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("error en scraping de Indeed: %s", result.Error)
	}

	return result.Offers, nil
}
