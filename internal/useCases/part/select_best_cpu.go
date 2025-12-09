package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
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

func (uc *SelectBestCPUUseCase) Execute(ctx context.Context, cpuPreference SelectBestCPUArgs) (entity.Part, error) {

	cpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      "CPU",
		Brand:         cpuPreference.BrandPreference,
		MaxPriceCents: cpuPreference.MaxPriceCents,
	})
	if err != nil {
		return entity.Part{}, err //TODO: adicionar erro custom depois
	}

	var bestCPU entity.Part

	for i, cpu := range cpus {
		if i == 0 || (cpu.Specs.PerformanceScore > bestCPU.Specs.PerformanceScore ||
			(cpu.Specs.PerformanceScore == bestCPU.Specs.PerformanceScore && cpu.PriceCents < bestCPU.PriceCents)) {
			bestCPU = *cpu
		}
	}

	return bestCPU, nil
}
