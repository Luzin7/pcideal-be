package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Luzin7/pcideal-be/internal/core/models"
)

type ScraperHTTPClient struct {
	BaseURL string
	Client  *http.Client
}

func NewScraperHTTPClient(baseURL string) *ScraperHTTPClient {
	return &ScraperHTTPClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ScraperHTTPClient) ScrapeAllCategories() ([]*models.Part, error) {
	response, err := s.Client.Get(fmt.Sprintf("%s/scrape-category", s.BaseURL))

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var parts []*models.Part

	err = json.NewDecoder(response.Body).Decode(&parts)

	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (s *ScraperHTTPClient) ScrapeProduct(productLink string) (*models.Part, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/scrape-product", s.BaseURL), nil)

	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("productLink", productLink)
	req.URL.RawQuery = query.Encode()

	resp, err := s.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var part models.Part

	err = json.NewDecoder(resp.Body).Decode(&part)

	if err != nil {
		return nil, err
	}

	return &part, nil
}
