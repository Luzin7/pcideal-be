package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type SelectBestGPUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestGPUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestGPUUseCase {
	return &SelectBestGPUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestGPUArgs struct {
	gpus []*entity.Part
}

func (uc *SelectBestGPUUseCase) Execute(ctx context.Context, args SelectBestGPUArgs) (entity.Part, *errors.ErrService) {
	var bestGPU entity.Part
	var partsToUpdate []dto.ProductLinkToUpdate

	for i, gpu := range args.gpus {
		if i == 0 {
			bestGPU = *gpu
			continue
		}

		if util.PartNeedToUpdate(gpu) {
			partsToUpdate = append(partsToUpdate, dto.ProductLinkToUpdate{
				ID:  gpu.ID.Hex(),
				Url: gpu.URL,
			})
		}

		if gpu.Specs.PerformanceScore > bestGPU.Specs.PerformanceScore {
			bestGPU = *gpu
			continue
		}

		if gpu.Specs.PerformanceScore < bestGPU.Specs.PerformanceScore {
			continue
		}

		if gpu.PriceCents < bestGPU.PriceCents {
			bestGPU = *gpu
		}
	}

	if len(partsToUpdate) > 0 {
		go func() {
			uc.UpdatePartsUseCase.Execute(context.Background(), partsToUpdate, "kabum")
		}()
	}

	return bestGPU, nil
}
