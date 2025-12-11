package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
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

	for i, psu := range args.psus {
		if i == 0 {
			bestPSU = *psu
			continue
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

	log.Printf("[SelectBestPSU] Selected: %s (Brand: %s, Price: %d, Wattage: %d, Efficiency: %d)", bestPSU.Model, bestPSU.Brand, bestPSU.PriceCents, bestPSU.Specs.Wattage, bestPSU.Specs.EfficiencyRating)

	return bestPSU, nil
}
