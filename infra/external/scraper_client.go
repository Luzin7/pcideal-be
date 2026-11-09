package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
)

const (
	defaultProductTimeout = 30 * time.Second
	defaultBulkTimeout    = 15 * time.Minute
	maxErrorBodyBytes     = 512
)

type ScraperStatusError struct {
	StatusCode int
	Body       string
}

func (e *ScraperStatusError) Error() string {
	return fmt.Sprintf("scraper returned status %d: %s", e.StatusCode, e.Body)
}

type ScraperHTTPClient struct {
	BaseURL        string
	ApiKey         string
	Client         *http.Client
	ProductTimeout time.Duration
	BulkTimeout    time.Duration
}

func NewScraperHTTPClient(baseURL string, apiKey string) *ScraperHTTPClient {
	return &ScraperHTTPClient{
		BaseURL:        baseURL,
		ApiKey:         apiKey,
		Client:         &http.Client{Timeout: defaultBulkTimeout},
		ProductTimeout: defaultProductTimeout,
		BulkTimeout:    defaultBulkTimeout,
	}
}

func (s *ScraperHTTPClient) ScrapeProduct(ctx context.Context, productLink string, storeName string) (*entity.Part, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/scrape-product", s.BaseURL, storeName), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.ApiKey)

	query := req.URL.Query()
	query.Add("productLink", productLink)
	req.URL.RawQuery = query.Encode()

	var part entity.Part
	if err := s.doJSON(ctx, req, s.ProductTimeout, &part); err != nil {
		return nil, err
	}

	return &part, nil
}

func (s *ScraperHTTPClient) UpdateProducts(ctx context.Context, links []dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error) {
	jsonBody, err := json.Marshal(links)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/update-products", s.BaseURL, storeName), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.ApiKey)

	var response dto.UpdateProductsResponse
	if err := s.doJSON(ctx, req, s.BulkTimeout, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (s *ScraperHTTPClient) doJSON(ctx context.Context, req *http.Request, timeout time.Duration, target any) error {
	requestCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := s.Client.Do(req.WithContext(requestCtx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return &ScraperStatusError{
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(body)),
		}
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
