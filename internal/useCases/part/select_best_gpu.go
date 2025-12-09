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
	BrandPreference string
	MaxPriceCents   int64
}

func (uc *SelectBestGPUUseCase) Execute(ctx context.Context, gpuPreference SelectBestGPUArgs) (entity.Part, error) {
	gpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "GPU",
		Brand:         gpuPreference.BrandPreference,
		MaxPriceCents: gpuPreference.MaxPriceCents,
	})
	if err != nil {
		return entity.Part{}, err //TODO: adicionar erro custom depois
	}

	var bestGPU entity.Part

	for i, gpu := range gpus {
		if i == 0 || (gpu.Specs.PerformanceScore > bestGPU.Specs.PerformanceScore ||
			(gpu.Specs.PerformanceScore == bestGPU.Specs.PerformanceScore && gpu.PriceCents < bestGPU.PriceCents)) {
			bestGPU = *gpu
		}
	}

	return bestGPU, nil
}
