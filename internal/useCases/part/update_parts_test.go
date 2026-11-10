package part

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Luzin7/pcideal-be/infra/external"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type failingScraperClient struct {
	err error
}

func (f failingScraperClient) ScrapeProduct(ctx context.Context, productLink string, storeName string) (*entity.Part, error) {
	return nil, f.err
}

func (f failingScraperClient) UpdateProducts(ctx context.Context, urls []dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error) {
	return nil, f.err
}

func captureLogOutput(t *testing.T, run func()) string {
	t.Helper()

	var buffer bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&buffer)
	defer log.SetOutput(previous)

	run()

	return buffer.String()
}

func TestUpdatePartsLogsContextOnScraperStatusError(t *testing.T) {
	urls := []dto.ProductLinkToUpdate{
		{ID: "part-a", Url: "https://store.example/a"},
		{ID: "part-b", Url: "https://store.example/b"},
	}
	client := failingScraperClient{err: &external.ScraperStatusError{StatusCode: 502, Body: "upstream error"}}
	uc := NewUpdatePartsUseCase(emptyPartRepository{}, client)

	var serviceErr *errors.ErrService
	output := captureLogOutput(t, func() {
		serviceErr = uc.Execute(context.Background(), urls, "pichau")
	})

	if serviceErr == nil || serviceErr.StatusCode != 502 || !strings.HasPrefix(serviceErr.Message, "Failed to scrape product") {
		t.Fatalf("expected scraping failed error, got %v", serviceErr)
	}

	if serviceErr.Message != "Failed to scrape products for store pichau (2 parts)" {
		t.Fatalf("expected batch error message to name store and batch size, got %q", serviceErr.Message)
	}

	for _, want := range []string{"pichau", "parts: 2", "part-a", "part-b"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected log output to contain %q, got %q", want, output)
		}
	}
}

func TestUpdatePartsLogsBoundedBatchContext(t *testing.T) {
	urls := make([]dto.ProductLinkToUpdate, 25)
	for i := range urls {
		urls[i] = dto.ProductLinkToUpdate{
			ID:  fmt.Sprintf("part-%02d", i+1),
			Url: fmt.Sprintf("https://store.example/%d", i+1),
		}
	}
	client := failingScraperClient{err: &external.ScraperStatusError{StatusCode: 502, Body: strings.Repeat("x", 4000)}}
	uc := NewUpdatePartsUseCase(emptyPartRepository{}, client)

	var serviceErr *errors.ErrService
	output := captureLogOutput(t, func() {
		serviceErr = uc.Execute(context.Background(), urls, "kabum")
	})

	if serviceErr == nil || serviceErr.StatusCode != 502 || serviceErr.Message != "Failed to scrape products for store kabum (25 parts)" {
		t.Fatalf("expected batch scraping failed error, got %v", serviceErr)
	}

	if got := strings.Count(output, "part-"); got != 20 {
		t.Fatalf("expected exactly 20 logged part IDs, got %d in %q", got, output)
	}

	for _, want := range []string{"kabum", "parts: 25", "part-20", "+5 more", "[truncated]"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected log output to contain %q, got %q", want, output)
		}
	}

	for _, tail := range []string{"part-21", "part-22", "part-23", "part-24", "part-25"} {
		if strings.Contains(output, tail) {
			t.Fatalf("expected log output not to contain %q, got %q", tail, output)
		}
	}
}

func TestUpdatePartsLogsContextOnTransportError(t *testing.T) {
	urls := []dto.ProductLinkToUpdate{
		{ID: "part-c", Url: "https://store.example/c"},
	}
	client := failingScraperClient{err: context.DeadlineExceeded}
	uc := NewUpdatePartsUseCase(emptyPartRepository{}, client)

	var serviceErr *errors.ErrService
	output := captureLogOutput(t, func() {
		serviceErr = uc.Execute(context.Background(), urls, "terabyte")
	})

	if serviceErr == nil || serviceErr.StatusCode != 502 || !strings.HasPrefix(serviceErr.Message, "Failed to scrape product") {
		t.Fatalf("expected scraping failed error, got %v", serviceErr)
	}

	for _, want := range []string{"terabyte", "parts: 1", "part-c", "context deadline exceeded"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected log output to contain %q, got %q", want, output)
		}
	}
}
