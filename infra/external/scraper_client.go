package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Luzin7/pcideal-be/infra/dto"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type ScraperHTTPClient struct {
	BaseURL string
	ApiKey  string
	Client  *http.Client
}

func NewScraperHTTPClient(baseURL string, apiKey string) *ScraperHTTPClient {
	return &ScraperHTTPClient{
		BaseURL: baseURL,
		ApiKey:  apiKey,
		Client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (s *ScraperHTTPClient) ScrapeAllCategories(storeName string) ([]*entity.Part, error) {
	response, err := s.Client.Get(fmt.Sprintf("%s/%s/scrape-category", s.BaseURL, storeName))

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var parts []*entity.Part

	err = json.NewDecoder(response.Body).Decode(&parts)

	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (s *ScraperHTTPClient) ScrapeProduct(productLink string, storeName string) (*entity.Part, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/scrape-product", s.BaseURL, storeName), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.ApiKey)

	query := req.URL.Query()
	query.Add("productLink", productLink)
	req.URL.RawQuery = query.Encode()

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Erro: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao chamar scraping API: %s", resp.Status)
	}

	var part entity.Part
	if err := json.NewDecoder(resp.Body).Decode(&part); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	log.Printf("Scraped part: %+v", part)
	return &part, nil
}

func (s *ScraperHTTPClient) UpdateProducts(links []*dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error) {
	jsonBody, err := json.Marshal(links)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/update-products", s.BaseURL, storeName), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", s.ApiKey)

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Erro: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("erro ao chamar scraping API: %s", resp.Status)
	}

	var partsWithID []*dto.PartWithID
	if err := json.NewDecoder(resp.Body).Decode(&partsWithID); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	return partsWithID, nil
}
