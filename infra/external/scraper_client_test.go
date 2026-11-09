package external

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/dto"
)

func TestScrapeProductReturnsNamedErrorOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("upstream down"))
	}))
	defer server.Close()

	client := NewScraperHTTPClient(server.URL, "key")

	_, err := client.ScrapeProduct(context.Background(), "https://example.com/product", "kabum")

	var statusErr *ScraperStatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected ScraperStatusError, got %v", err)
	}
	if statusErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, statusErr.StatusCode)
	}
}

func TestUpdateProductsReturnsNamedErrorOnNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewScraperHTTPClient(server.URL, "key")

	_, err := client.UpdateProducts(context.Background(), []dto.ProductLinkToUpdate{}, "kabum")

	var statusErr *ScraperStatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected ScraperStatusError, got %v", err)
	}
	if statusErr.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, statusErr.StatusCode)
	}
}

func TestScrapeProductHonoursPerCallTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewScraperHTTPClient(server.URL, "key")
	client.ProductTimeout = 5 * time.Millisecond

	_, err := client.ScrapeProduct(context.Background(), "https://example.com/product", "kabum")

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded, got %v", err)
	}
}
