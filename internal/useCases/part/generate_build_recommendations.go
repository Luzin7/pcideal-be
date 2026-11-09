package part

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type GenerateBuildRecommendationsUseCase struct {
	partRepository   repository.PartRepository
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
		updatePartsUC:    updatePartsUC,
		selectBestCPUUC:  selectBestCPUUC,
		selectBestGPUUC:  selectBestGPUUC,
		selectBestPSUUC:  selectBestPSUUC,
		selectBestRAMUC:  selectBestRAMUC,
		selectBestMOBOUC: selectBestMOBOUC,
		selectBestSSDUC:  selectBestSSDUC,
	}
}

type buildTierInput struct {
	args        dto.GenerateBuildRecommendationsDTO
	allocations map[entity.PartType]float64
	buildType   string
	multiplier  float64
}

func (uc *GenerateBuildRecommendationsUseCase) Execute(ctx context.Context, args dto.GenerateBuildRecommendationsDTO) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	budgetCents := args.BudgetCents
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
	var firstBuildErr *errors.ErrService

	for _, config := range buildConfigs {
		build, buildErr := uc.buildTier(ctx, buildTierInput{
			args:        args,
			allocations: allocations,
			buildType:   config.buildType,
			multiplier:  config.multiplier,
		})
		if buildErr != nil {
			log.Printf("[GenerateBuildRecommendations] skipping %s build: %v", config.buildType, buildErr)
			if firstBuildErr == nil {
				firstBuildErr = buildErr
			}
			continue
		}

		recommendedBuilds = append(recommendedBuilds, *build)
	}

	if len(recommendedBuilds) == 0 {
		return nil, firstBuildErr
	}

	return &presenters.RecommendedBuildsPresenter{
		Builds: recommendedBuilds,
	}, nil
}

func (uc *GenerateBuildRecommendationsUseCase) buildTier(ctx context.Context, input buildTierInput) (*presenters.RecommendationBuild, *errors.ErrService) {
	buildBudgetCents := int64(float64(input.args.BudgetCents) * input.multiplier)

	maxCpuBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypeCPU])

	cpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:           entity.TypeCPU,
		Brand:              input.args.CpuPreference,
		MaxPriceCents:      maxCpuBudgetCents,
		IntegratedGraphics: input.args.UsageType == "WORK",
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
	if err := ensurePartSelected(selectedCpu, "CPU"); err != nil {
		return nil, err
	}

	maxMoboBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypeMobo])

	var moboBrandPreference string
	if input.args.CpuPreference != "" && input.args.CpuPreference != "no-preference" {
		moboBrandPreference = input.args.CpuPreference
	} else {
		moboBrandPreference = selectedCpu.Brand
	}

	mobos, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypeMobo,
		Brand:         moboBrandPreference,
		Socket:        selectedCpu.Specs.Socket,
		MaxPriceCents: maxMoboBudgetCents,
		MemoryType:    selectedCpu.Specs.MemoryType,
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
	if err := ensurePartSelected(selectedMobo, "MOTHERBOARD"); err != nil {
		return nil, err
	}

	maxGpuBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypeGPU])

	var selectedGpu entity.Part

	selectedGpu = entity.Part{
		Brand:         "Integrado",
		Model:         "Processador",
		AffiliatedURL: selectedCpu.AffiliatedURL,
		PriceCents:    0,
	}

	if input.args.UsageType != "WORK" {
		gpus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
			PartType:      entity.TypeGPU,
			Brand:         input.args.GpuPreference,
			MaxPriceCents: maxGpuBudgetCents,
		})
		if err != nil {
			log.Printf("[SelectBestGPU] Error querying database: %v", err)
			return nil, errors.ErrBuildAttemptNotFound()
		}

		selecteBestdGpu, errGpu := uc.selectBestGPUUC.Execute(ctx, SelectBestGPUArgs{
			gpus: gpus,
		})
		if errGpu != nil {
			return nil, errors.ErrBuildAttemptNotFound()
		}
		if err := ensurePartSelected(selecteBestdGpu, "GPU"); err != nil {
			return nil, err
		}

		selectedGpu = selecteBestdGpu
	}

	maxPsuBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypePSU])

	psus, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypePSU,
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
	if err := ensurePartSelected(selectedPsu, "PSU"); err != nil {
		return nil, err
	}

	maxRamBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypeRAM])

	rams, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypeRAM,
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
	if err := ensurePartSelected(selectedRam, "RAM"); err != nil {
		return nil, err
	}

	maxSsdBudgetCents := int64(float64(buildBudgetCents) * input.allocations[entity.TypeSSD])

	ssds, err := uc.partRepository.FindPartByTypeAndBrandWithMaxPrice(ctx, repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypeSSD,
		MaxPriceCents: maxSsdBudgetCents,
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
	if err := ensurePartSelected(selectedSSD, "SSD"); err != nil {
		return nil, err
	}

	return &presenters.RecommendationBuild{
		BuildType: input.buildType,
		Parts: presenters.BuildParts{
			CPU:            &selectedCpu,
			Motherboard:    &selectedMobo,
			GPU:            &selectedGpu,
			PSU:            &selectedPsu,
			PrimaryStorage: &selectedSSD,
			RAM:            &selectedRam,
		},
		Budget:     input.args.BudgetCents,
		BuildValue: selectedCpu.PriceCents + selectedGpu.PriceCents + selectedMobo.PriceCents + selectedPsu.PriceCents + selectedRam.PriceCents + selectedSSD.PriceCents,
	}, nil
}

func ensurePartSelected(part entity.Part, partType string) *errors.ErrService {
	if part.Model != "" {
		return nil
	}

	return errors.ErrNoCompatibleParts(partType)
}
