package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestRAMUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestRAMUseCase(partRepository repository.PartRepository) *SelectBestRAMUseCase {
	return &SelectBestRAMUseCase{
		partRepository: partRepository,
	}
}

type SelectBestRAMArgs struct {
	BrandPreference string
	MaxPriceCents   int64
}

func (uc *SelectBestRAMUseCase) Execute(ctx context.Context, ramPreference SelectBestRAMArgs) (entity.Part, *errors.ErrService) {
	log.Printf("[SelectBestRAM] Filtering - Brand: %s, MaxPrice: %d", ramPreference.BrandPreference, ramPreference.MaxPriceCents)

	rams, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "RAM",
		Brand:         ramPreference.BrandPreference,
		MaxPriceCents: ramPreference.MaxPriceCents,
	})
	if err != nil {
		log.Printf("[SelectBestRAM] Error querying database: %v", err)
		return entity.Part{}, errors.New("Failed to select best RAM", 500)
	}

	log.Printf("[SelectBestRAM] Found %d RAMs", len(rams))

	if len(rams) == 0 {
		log.Printf("[SelectBestRAM] No RAMs found matching criteria")
		return entity.Part{}, errors.New("No RAM found matching the criteria", 404)
	}

	var bestRAM entity.Part

	for i, ram := range rams {
		if ram.Specs.CasLatency <= 0 {
			continue
		}
		if i == 0 {
			bestRAM = *ram
			continue
		}

		if ram.Specs.PerformanceScore > bestRAM.Specs.PerformanceScore {
			bestRAM = *ram
			continue
		}

		if ram.Specs.PerformanceScore < bestRAM.Specs.PerformanceScore {
			continue
		}

		if ram.Specs.CasLatency < bestRAM.Specs.CasLatency {
			bestRAM = *ram
			continue
		}

		if ram.Specs.CasLatency > bestRAM.Specs.CasLatency {
			continue
		}

		if ram.PriceCents < bestRAM.PriceCents {
			bestRAM = *ram
		}
	}

	log.Printf("[SelectBestRAM] Selected: %s (Brand: %s, Price: %d, Score: %d, CAS: %d)", bestRAM.Model, bestRAM.Brand, bestRAM.PriceCents, bestRAM.Specs.PerformanceScore, bestRAM.Specs.CasLatency)

	return bestRAM, nil
}
