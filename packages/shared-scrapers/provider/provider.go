package provider

import "context"

type JobOffer struct {
	ID             string   `json:"id"`
	Platform       string   `json:"platform"` // "infojobs", "indeed", "linkedin"
	Title          string   `json:"title"`
	Company        string   `json:"company"`
	Location       string   `json:"location"`
	RawDescription string   `json:"raw_description"`
	URL            string   `json:"url"`
	Salary         string   `json:"salary,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

type JobProvider interface {
	Search(ctx context.Context, query string, location string) ([]JobOffer, error)
}
