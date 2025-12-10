package part

import (
	"context"
	"log"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type UpdatePartUseCase struct {
	partRepository repository.PartRepository
	scraperClient  repository.ScraperClientRepository
}

func NewUpdatePartUseCase(partRepository repository.PartRepository, scraperClient repository.ScraperClientRepository) *UpdatePartUseCase {
	return &UpdatePartUseCase{
		partRepository: partRepository,
		scraperClient:  scraperClient,
	}
}

func (uc *UpdatePartUseCase) Execute(partId string) *errors.ErrService {
	ctx := context.Background()
	part, err := uc.partRepository.GetPartByID(ctx, partId)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	if part == nil {
		return errors.ErrNotFound("part")
	}

	scrapedPart, err := uc.scraperClient.ScrapeProduct(ctx, part.URL, "kabum")

	if err != nil {
		log.Printf("Error scraping product for part %s: %v", partId, err)
		return errors.ErrScrapingFailed(partId)
	}

	part.Brand = scrapedPart.Brand
	part.Model = scrapedPart.Model
	part.Specs = scrapedPart.Specs
	part.PriceCents = scrapedPart.PriceCents
	part.UpdatedAt = time.Now()

	err = uc.partRepository.UpdatePart(ctx, partId, part)

	if err != nil {
		log.Printf("Error updating part: %v", err)
		return errors.ErrInternalServerError()
	}

	return nil
}
