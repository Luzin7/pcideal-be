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
	cpus []*entity.Part
}

func (uc *SelectBestCPUUseCase) Execute(ctx context.Context, args SelectBestCPUArgs) (entity.Part, *errors.ErrService) {
	var bestCPU entity.Part

	for i, cpu := range args.cpus {
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
