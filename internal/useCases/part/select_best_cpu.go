package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestCPUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestCPUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestCPUUseCase {
	return &SelectBestCPUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestCPUArgs struct {
	cpus []*entity.Part
}

func (uc *SelectBestCPUUseCase) Execute(ctx context.Context, args SelectBestCPUArgs) (entity.Part, *errors.ErrService) {
	var bestCPU entity.Part
	var partsToUpdate []dto.ProductLinkToUpdate

	for i, cpu := range args.cpus {
		if i == 0 {
			bestCPU = *cpu
			continue
		}

		if util.PartNeedToUpdate(cpu) {
			partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
				ID:  cpu.ID.Hex(),
				Url: cpu.URL,
			})
		}

		if cpu.Specs.PerformanceScore > bestCPU.Specs.PerformanceScore {
			bestCPU = *cpu
			continue
		}

		if cpu.Specs.PerformanceScore < bestCPU.Specs.PerformanceScore {
			continue
		}

		if cpu.PriceCents < bestCPU.PriceCents {
			bestCPU = *cpu
		}
	}

	if len(partsToUpdate) > 0 {
		go func() {
			uc.UpdatePartsUseCase.Execute(context.Background(), partsToUpdate, "kabum")
		}()
	}

	return bestCPU, nil
}
