package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestGPUUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestGPUUseCase(partRepository repository.PartRepository) *SelectBestGPUUseCase {
	return &SelectBestGPUUseCase{
		partRepository: partRepository,
	}
}

type SelectBestGPUArgs struct {
	gpus []*entity.Part
}

func (uc *SelectBestGPUUseCase) Execute(ctx context.Context, args SelectBestGPUArgs) (entity.Part, *errors.ErrService) {
	var bestGPU entity.Part

	for i, gpu := range args.gpus {
		if i == 0 {
			bestGPU = *gpu
			continue
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

	log.Printf("[SelectBestGPU] Selected: %s (Brand: %s, Price: %d, Score: %d, MinPSU: %d)", bestGPU.Model, bestGPU.Brand, bestGPU.PriceCents, bestGPU.Specs.PerformanceScore, bestGPU.Specs.MinPSUWatts)

	return bestGPU, nil
}
