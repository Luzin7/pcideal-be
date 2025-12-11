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
	mobos []*entity.Part
}

func (uc *SelectBestMOBOUseCase) Execute(ctx context.Context, args SelectBestMOBOArgs) (entity.Part, *errors.ErrService) {
	var bestMOBO entity.Part

	for i, mobo := range args.mobos {
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
