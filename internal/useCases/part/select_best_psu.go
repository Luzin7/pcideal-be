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
	BrandPreference string
	MaxPriceCents   int64
	MinPSUWatts     int16
}

func (uc *SelectBestPSUUseCase) Execute(ctx context.Context, psuPreference SelectBestPSUArgs) (entity.Part, *errors.ErrService) {
	log.Printf("[SelectBestPSU] Filtering - Brand: %s, MaxPrice: %d, MinWatts: %d", psuPreference.BrandPreference, psuPreference.MaxPriceCents, psuPreference.MinPSUWatts)

	psus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "PSU",
		Brand:         psuPreference.BrandPreference,
		MaxPriceCents: psuPreference.MaxPriceCents,
		MinPSUWatts:   psuPreference.MinPSUWatts,
	})
	if err != nil {
		log.Printf("[SelectBestPSU] Error querying database: %v", err)
		return entity.Part{}, errors.New("Failed to select best PSU", 500)
	}

	log.Printf("[SelectBestPSU] Found %d PSUs", len(psus))

	if len(psus) == 0 {
		log.Printf("[SelectBestPSU] No PSUs found matching criteria")
		return entity.Part{}, errors.New("No PSU found matching the criteria", 404)
	}

	var bestPSU entity.Part

	for i, psu := range psus {
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
