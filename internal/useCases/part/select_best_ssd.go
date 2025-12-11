package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestSSDUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestSSDUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestSSDUseCase {
	return &SelectBestSSDUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestSSDArgs struct {
	ssds []*entity.Part
}

func (uc *SelectBestSSDUseCase) Execute(ctx context.Context, args SelectBestSSDArgs) (entity.Part, *errors.ErrService) {
	var bestSSD entity.Part
	var partsToUpdate []dto.ProductLinkToUpdate

	for i, ssd := range args.ssds {
		if i == 0 {
			bestSSD = *ssd
			continue
		}

		if util.PartNeedToUpdate(ssd) {
			partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
				ID:  ssd.ID.Hex(),
				Url: ssd.URL,
			})
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

	if len(partsToUpdate) > 0 {
		go func() {
			uc.UpdatePartsUseCase.Execute(context.Background(), partsToUpdate, "kabum")
		}()
	}

	return bestSSD, nil
}
