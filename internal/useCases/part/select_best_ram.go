package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestRAMUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestRAMUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestRAMUseCase {
	return &SelectBestRAMUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestRAMArgs struct {
	rams []*entity.Part
}

func (uc *SelectBestRAMUseCase) Execute(ctx context.Context, args SelectBestRAMArgs) (entity.Part, *errors.ErrService) {
	var bestRAM entity.Part
	var partsToUpdate []dto.ProductLinkToUpdate

	for i, ram := range args.rams {
		if ram.Specs.CasLatency <= 0 {
			continue
		}
		if i == 0 {
			bestRAM = *ram
			continue
		}

		if util.PartNeedToUpdate(ram) {
			partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
				ID:  ram.ID.Hex(),
				Url: ram.URL,
			})
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

	if len(partsToUpdate) > 0 {
		go func() {
			uc.UpdatePartsUseCase.Execute(context.Background(), partsToUpdate, "kabum")
		}()
	}

	return bestRAM, nil
}
