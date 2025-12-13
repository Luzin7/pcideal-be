package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type GenerateBuildRecommendationsUseCase struct {
	partRepository   repository.PartRepository
	scraperClient    repository.ScraperClientRepository
	googleAIClient   repository.GoogleAIRepository
	updatePartsUC    *UpdatePartsUseCase
	selectBestCPUUC  *SelectBestCPUUseCase
	selectBestGPUUC  *SelectBestGPUUseCase
	selectBestPSUUC  *SelectBestPSUUseCase
	selectBestRAMUC  *SelectBestRAMUseCase
	selectBestMOBOUC *SelectBestMOBOUseCase
	selectBestSSDUC  *SelectBestSSDUseCase
}

func NewGenerateBuildRecommendationsUseCase(
	partRepository repository.PartRepository,
	scraperClient repository.ScraperClientRepository,
	googleAIClient repository.GoogleAIRepository,
	updatePartsUC *UpdatePartsUseCase,
	selectBestCPUUC *SelectBestCPUUseCase,
	selectBestGPUUC *SelectBestGPUUseCase,
	selectBestPSUUC *SelectBestPSUUseCase,
	selectBestRAMUC *SelectBestRAMUseCase,
	selectBestMOBOUC *SelectBestMOBOUseCase,
	selectBestSSDUC *SelectBestSSDUseCase,
) *GenerateBuildRecommendationsUseCase {
	return &GenerateBuildRecommendationsUseCase{
		partRepository:   partRepository,
		scraperClient:    scraperClient,
		googleAIClient:   googleAIClient,
		updatePartsUC:    updatePartsUC,
		selectBestCPUUC:  selectBestCPUUC,
		selectBestGPUUC:  selectBestGPUUC,
		selectBestPSUUC:  selectBestPSUUC,
		selectBestRAMUC:  selectBestRAMUC,
		selectBestMOBOUC: selectBestMOBOUC,
		selectBestSSDUC:  selectBestSSDUC,
	}
}

func (uc *GenerateBuildRecommendationsUseCase) Execute(ctx context.Context, args dto.GenerateBuildRecommendationsDTO) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	budgetCents := util.ConvertToCents(args.Budget)
	strategy := entity.GetStrategy(args.UsageType, budgetCents)
	allocations := strategy.GetAllocations()

	buildConfigs := []struct {
		buildType  string
		multiplier float64
	}{
		{"ECONOMIC", 0.85},
		{"BALANCED", 1.0},
		{"PERFORMANCE", 1.15},
	}

	var recommendedBuilds []presenters.RecommendationBuild

	for _, config := range buildConfigs {
		var finalBuild *presenters.RecommendationBuild

		buildBudgetCents := int64(float64(budgetCents) * config.multiplier)

		maxCpuBudgetCents := int64(float64(buildBudgetCents) * allocations[entity.TypeCPU])

		cpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "CPU",
			Brand:         args.CpuPreference,
			MaxPriceCents: maxCpuBudgetCents,
		})
		if err != nil {
			log.Printf("[SelectBestCPU] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedCpu, errCpu := uc.selectBestCPUUC.Execute(ctx, SelectBestCPUArgs{
			cpus: cpus,
		})
		if errCpu != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		maxMoboBudgetCents := int64(float64(buildBudgetCents) * allocations[entity.TypeMobo])

		mobos, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "MOTHERBOARD",
			Brand:         args.CpuPreference,
			Socket:        selectedCpu.Specs.Socket,
			MaxPriceCents: maxMoboBudgetCents,
		})
		if err != nil {
			log.Printf("[SelectBestMOBO] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedMobo, errMobo := uc.selectBestMOBOUC.Execute(ctx, SelectBestMOBOArgs{
			mobos: mobos,
		})
		if errMobo != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		maxGpuBudgetCents := int64(float64(buildBudgetCents) * allocations[entity.TypeGPU])

		gpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "GPU",
			Brand:         args.GpuPreference,
			MaxPriceCents: maxGpuBudgetCents,
		})
		if err != nil {
			log.Printf("[SelectBestGPU] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedGpu, errGpu := uc.selectBestGPUUC.Execute(ctx, SelectBestGPUArgs{
			gpus: gpus,
		})
		if errGpu != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		maxPsuBudgetCents := int64(float64(buildBudgetCents) * allocations[entity.TypePSU])

		psus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "PSU",
			MaxPriceCents: maxPsuBudgetCents,
			MinPSUWatts:   selectedCpu.Specs.MinPSUWatts + 100,
		})
		if err != nil {
			log.Printf("[SelectBestPSU] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedPsu, errPsu := uc.selectBestPSUUC.Execute(ctx, SelectBestPSUArgs{
			psus: psus,
		})
		if errPsu != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		maxRamBudgetCents := int64(float64(buildBudgetCents) * allocations[entity.TypeRAM])

		rams, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "RAM",
			MaxPriceCents: maxRamBudgetCents,
			MemoryType:    selectedCpu.Specs.MemoryType,
		})
		if err != nil {
			log.Printf("[SelectBestRAM] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedRam, errRam := uc.selectBestRAMUC.Execute(ctx, SelectBestRAMArgs{
			rams: rams,
		})
		if errRam != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		ssds, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      "SSD",
			MaxPriceCents: maxRamBudgetCents,
		})
		if err != nil {
			log.Printf("[SelectBestSSD] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selectedSSD, errSSD := uc.selectBestSSDUC.Execute(ctx, SelectBestSSDArgs{
			ssds: ssds,
		})
		if errSSD != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}

		finalBuild = &presenters.RecommendationBuild{
			BuildType: config.buildType,
			Parts: presenters.BuildParts{
				CPU:            &selectedCpu,
				Motherboard:    &selectedMobo,
				GPU:            &selectedGpu,
				PSU:            &selectedPsu,
				PrimaryStorage: &selectedSSD,
				RAM:            &selectedRam,
			},
			Budget:     args.Budget,
			BuildValue: selectedCpu.PriceCents + selectedGpu.PriceCents + selectedMobo.PriceCents + selectedPsu.PriceCents + selectedRam.PriceCents + selectedSSD.PriceCents,
		}

		finalBuild.Summary = "Análise indisponível no momento. Em breve você terá uma análise detalhada do build aqui."

		recommendedBuilds = append(recommendedBuilds, *finalBuild)

	}

	return &presenters.RecommendedBuildsPresenter{
		Builds: recommendedBuilds,
	}, nil
}
