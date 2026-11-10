package part

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type UpdatePartsUseCase struct {
	partRepository repository.PartRepository
	scraperClient  repository.ScraperClientRepository
}

func NewUpdatePartsUseCase(partRepository repository.PartRepository, scraperClient repository.ScraperClientRepository) *UpdatePartsUseCase {
	return &UpdatePartsUseCase{
		partRepository: partRepository,
		scraperClient:  scraperClient,
	}
}

const (
	maxLoggedPartIDs     = 20
	maxLoggedErrorLength = 512
)

func partIDs(urls []dto.ProductLinkToUpdate) []string {
	ids := make([]string, len(urls))
	for i, url := range urls {
		ids[i] = url.ID
	}

	return ids
}

func loggedPartIDs(urls []dto.ProductLinkToUpdate) string {
	ids := partIDs(urls)
	if len(ids) <= maxLoggedPartIDs {
		return strings.Join(ids, ", ")
	}

	return fmt.Sprintf("%s, +%d more", strings.Join(ids[:maxLoggedPartIDs], ", "), len(ids)-maxLoggedPartIDs)
}

func loggedError(err error) string {
	message := []rune(err.Error())
	if len(message) <= maxLoggedErrorLength {
		return string(message)
	}

	return fmt.Sprintf("%s [truncated]", string(message[:maxLoggedErrorLength]))
}

func (uc *UpdatePartsUseCase) Execute(ctx context.Context, urls []dto.ProductLinkToUpdate, store string) *errors.ErrService {
	if len(urls) <= 0 {
		return errors.New("urls cannot be empty", 400)
	}

	updatedParts, err := uc.scraperClient.UpdateProducts(ctx, urls, store)

	if err != nil {
		log.Printf("Error scraping products for store %s (parts: %d, ids: %s): %s", store, len(urls), loggedPartIDs(urls), loggedError(err))
		return errors.ErrScrapingFailedBatch(store, len(urls))
	}

	for _, part := range updatedParts {
		if part.Part == nil {
			log.Printf("Warning: part.Part is nil for ID %s, skipping update", part.ID)
			continue
		}

		part.Part.UpdatedAt = time.Now().UTC()

		err = uc.partRepository.UpdatePart(ctx, part.ID, part.Part)
		if err != nil {
			log.Printf("Error updating part: %v", err)
			return errors.ErrInternalServerError()
		}
	}

	return nil
}
