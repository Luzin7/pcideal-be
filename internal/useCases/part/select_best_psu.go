package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestPSUUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestPSUUseCase(partRepository repository.PartRepository) *SelectBestPSUUseCase {
	return &SelectBestPSUUseCase{
		partRepository: partRepository,
	}
}

type SelectBestPSUArgs struct {
	psus []*entity.Part
}

func (uc *SelectBestPSUUseCase) Execute(ctx context.Context, args SelectBestPSUArgs) (entity.Part, *errors.ErrService) {
	var bestPSU entity.Part
	var partsToUpdate []dto.ProductLinkToUpdate

	for i, psu := range args.psus {
		if i == 0 {
			bestPSU = *psu
			continue
		}

		if util.PartNeedToUpdate(psu) {
			partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
				ID:  psu.ID.Hex(),
				Url: psu.URL,
			})
		}

		if psu.Specs.EfficiencyRating > bestPSU.Specs.EfficiencyRating {
			bestPSU = *psu
			continue
		}

		if psu.Specs.EfficiencyRating < bestPSU.Specs.EfficiencyRating {
			continue
		}

		if psu.PriceCents < bestPSU.PriceCents {
			bestPSU = *psu
		}
	}

	if len(partsToUpdate) > 0 {
		go func() {
			uc.partRepository.UpdateParts(context.Background(), partsToUpdate, "kabum")
		}()
	}

	log.Printf("[SelectBestPSU] Selected: %s (Brand: %s, Price: %d, Wattage: %d, Efficiency: %d)", bestPSU.Model, bestPSU.Brand, bestPSU.PriceCents, bestPSU.Specs.Wattage, bestPSU.Specs.EfficiencyRating)

	return bestPSU, nil
}
