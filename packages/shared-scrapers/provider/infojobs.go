package provider

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type InfoJobsProvider struct {
	Client *http.Client
}

func NewInfoJobsProvider() *InfoJobsProvider {
	return &InfoJobsProvider{
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (p *InfoJobsProvider) Search(ctx context.Context, query, location string) ([]JobOffer, error) {
	searchURL := fmt.Sprintf(
		"https://www.infojobs.net/jobsearch/search-results/list.xhtml?keyword=%s",
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición a InfoJobs: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando petición a InfoJobs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("InfoJobs respondió con estado HTTP: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error parseando HTML de InfoJobs: %w", err)
	}

	var offers []JobOffer

	// Selectores actualizados para la lista de ofertas de InfoJobs
	doc.Find("div.ij-OfferCardContent, li[class*='OfferCard'], article").Each(func(i int, s *goquery.Selection) {
		titleElem := s.Find("a[class*='title'], a.ij-OfferCardContent-description-title-link")
		title := strings.TrimSpace(titleElem.Text())
		link, _ := titleElem.Attr("href")

		company := strings.TrimSpace(s.Find("a[class*='subtitle'], a.ij-OfferCardContent-description-subtitle-link").Text())
		loc := strings.TrimSpace(s.Find("span[class*='property'], span.ij-OfferCardContent-description-list-item-truncate").First().Text())

		if link != "" && !strings.HasPrefix(link, "http") {
			link = "https:" + link
		}

		if title != "" && link != "" {
			offers = append(offers, JobOffer{
				Platform: "infojobs",
				Title:    title,
				Company:  company,
				Location: loc,
				URL:      link,
			})
		}
	})

	return offers, nil
}
