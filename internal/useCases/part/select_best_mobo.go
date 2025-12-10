package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestMOBOUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestMOBOUseCase(partRepository repository.PartRepository) *SelectBestMOBOUseCase {
	return &SelectBestMOBOUseCase{
		partRepository: partRepository,
	}
}

type SelectBestMOBOArgs struct {
	Brand         string
	Socket        string
	MaxPriceCents int64
}

func (uc *SelectBestMOBOUseCase) Execute(ctx context.Context, moboPreference SelectBestMOBOArgs) (entity.Part, *errors.ErrService) {
	log.Printf("[SelectBestMOBO] Filtering - Socket: %s, MaxPrice: %d", moboPreference.Socket, moboPreference.MaxPriceCents)

	mobos, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "MOTHERBOARD",
		Brand:         moboPreference.Brand,
		Socket:        moboPreference.Socket,
		MaxPriceCents: moboPreference.MaxPriceCents,
	})
	if err != nil {
		log.Printf("[SelectBestMOBO] Error querying database: %v", err)
		return entity.Part{}, errors.New("Failed to select best MOBO", 500)
	}

	log.Printf("[SelectBestMOBO] Found %d MOBOs", len(mobos))

	if len(mobos) == 0 {
		log.Printf("[SelectBestMOBO] No MOBOs found matching criteria")
		return entity.Part{}, errors.New("No MOBO found matching the criteria", 404)
	}

	var bestMOBO entity.Part

	for i, mobo := range mobos {
		if i == 0 {
			bestMOBO = *mobo
			continue
		}

		if mobo.Specs.TierScore > bestMOBO.Specs.TierScore {
			bestMOBO = *mobo
			continue
		}

		if mobo.Specs.TierScore < bestMOBO.Specs.TierScore {
			continue
		}

		if mobo.PriceCents < bestMOBO.PriceCents {
			bestMOBO = *mobo
		}
	}

	log.Printf("[SelectBestMOBO] Selected: %s (Brand: %s, Price: %d, Socket: %s, TierScore: %d)", bestMOBO.Model, bestMOBO.Brand, bestMOBO.PriceCents, bestMOBO.Specs.Socket, bestMOBO.Specs.TierScore)

	return bestMOBO, nil
}
