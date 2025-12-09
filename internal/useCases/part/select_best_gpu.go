package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
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
	brandPreference string
	maxPriceCents   int64
}

func (uc *SelectBestGPUUseCase) Execute(gpuPreference SelectBestGPUArgs) (entity.Part, error) {
	ctx := context.TODO()

	gpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "GPU",
		Brand:         gpuPreference.brandPreference,
		MaxPriceCents: gpuPreference.maxPriceCents,
	})
	if err != nil {
		return entity.Part{}, err //TODO: adicionar erro custom depois
	}

	var bestGPU entity.Part

	for i, gpu := range gpus {
		if i == 0 || gpu.Specs.PerformanceScore >= bestGPU.Specs.PerformanceScore && gpu.PriceCents <= gpuPreference.maxPriceCents {
			bestGPU = *gpu
		}
	}

	return bestGPU, nil
}
