package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestPSUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestPSUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestPSUUseCase {
	return &SelectBestPSUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
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
			uc.UpdatePartsUseCase.Execute(context.Background(), partsToUpdate, "kabum")
		}()
	}

	return bestPSU, nil
}
