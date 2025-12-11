package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestSSDUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestSSDUseCase(partRepository repository.PartRepository) *SelectBestSSDUseCase {
	return &SelectBestSSDUseCase{
		partRepository: partRepository,
	}
}

type SelectBestSSDArgs struct {
	ssds []*entity.Part
}

func (uc *SelectBestSSDUseCase) Execute(ctx context.Context, args SelectBestSSDArgs) (entity.Part, *errors.ErrService) {
	var bestSSD entity.Part

	for i, ssd := range args.ssds {
		if i == 0 {
			bestSSD = *ssd
			continue
		}

		if ssd.Specs.EfficiencyRating > bestSSD.Specs.EfficiencyRating {
			bestSSD = *ssd
			continue
		}

		if ssd.Specs.EfficiencyRating < bestSSD.Specs.EfficiencyRating {
			continue
		}

		if ssd.PriceCents < bestSSD.PriceCents {
			bestSSD = *ssd
		}
	}

	log.Printf("[SelectBestSSD] Selected: %s (Brand: %s, Price: %d, Wattage: %d, Efficiency: %d)", bestSSD.Model, bestSSD.Brand, bestSSD.PriceCents, bestSSD.Specs.Wattage, bestSSD.Specs.EfficiencyRating)

	return bestSSD, nil
}
