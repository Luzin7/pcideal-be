package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestCPUUseCase struct {
	partRepository repository.PartRepository
}

func NewSelectBestCPUUseCase(partRepository repository.PartRepository) *SelectBestCPUUseCase {
	return &SelectBestCPUUseCase{
		partRepository: partRepository,
	}
}

type SelectBestCPUArgs struct {
	BrandPreference string
	MaxPriceCents   int64
}

func (uc *SelectBestCPUUseCase) Execute(ctx context.Context, cpuPreference SelectBestCPUArgs) (entity.Part, *errors.ErrService) {
	log.Printf("[SelectBestCPU] Filtering - Brand: %s, MaxPrice: %d", cpuPreference.BrandPreference, cpuPreference.MaxPriceCents)

	cpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "CPU",
		Brand:         cpuPreference.BrandPreference,
		MaxPriceCents: cpuPreference.MaxPriceCents,
	})
	if err != nil {
		log.Printf("[SelectBestCPU] Error querying database: %v", err)
		return entity.Part{}, errors.New("Failed to select best CPU", 500)
	}

	log.Printf("[SelectBestCPU] Found %d CPUs", len(cpus))

	if len(cpus) == 0 {
		log.Printf("[SelectBestCPU] No CPUs found matching criteria")
		return entity.Part{}, errors.New("No CPU found matching the criteria", 404)
	}

	var bestCPU entity.Part

	for i, cpu := range cpus {
		if i == 0 {
			bestCPU = *cpu
			continue
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

	log.Printf("[SelectBestCPU] Selected: %s (Brand: %s, Price: %d, Score: %d, Socket: %s)", bestCPU.Model, bestCPU.Brand, bestCPU.PriceCents, bestCPU.Specs.PerformanceScore, bestCPU.Specs.Socket)

	return bestCPU, nil
}
