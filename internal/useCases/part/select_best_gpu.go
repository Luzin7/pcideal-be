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
	BrandPreference string
	MaxPriceCents   int64
}

func (uc *SelectBestGPUUseCase) Execute(ctx context.Context, gpuPreference SelectBestGPUArgs) (entity.Part, *errors.ErrService) {
	log.Printf("[SelectBestGPU] Filtering - Brand: %s, MaxPrice: %d", gpuPreference.BrandPreference, gpuPreference.MaxPriceCents)

	gpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "GPU",
		Brand:         gpuPreference.BrandPreference,
		MaxPriceCents: gpuPreference.MaxPriceCents,
	})
	if err != nil {
		log.Printf("[SelectBestGPU] Error querying database: %v", err)
		return entity.Part{}, errors.New("Failed to select best GPU", 500)
	}

	log.Printf("[SelectBestGPU] Found %d GPUs", len(gpus))

	if len(gpus) == 0 {
		log.Printf("[SelectBestGPU] No GPUs found matching criteria")
		return entity.Part{}, errors.New("No GPU found matching the criteria", 404)
	}

	var bestGPU entity.Part

	for i, gpu := range gpus {
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
