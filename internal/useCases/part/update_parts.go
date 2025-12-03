package part

import (
	"context"
	"log"
	"time"

	"github.com/Luzin7/pcideal-be/infra/dto"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
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

func (uc *UpdatePartsUseCase) Execute(urls []*dto.ProductLinkToUpdate) *errors.ErrService {
	ctx := context.Background()
	if len(urls) <= 0 {
		return errors.New("urls cannot be empty", 400)
	}

	updatedParts, err := uc.scraperClient.UpdateProducts(ctx, urls, "kabum")

	if err != nil {
		log.Printf("Error scraping product for part %v", err)
		return errors.ErrScrapingFailed("urls")
	}

	for _, part := range updatedParts {
		part.Part.UpdatedAt = time.Now()

		err = uc.partRepository.UpdatePart(ctx, part.ID, part.Part)
		if err != nil {
			log.Printf("Error updating part: %v", err)
			return errors.ErrInternalServerError()
		}
	}

	return nil
}
