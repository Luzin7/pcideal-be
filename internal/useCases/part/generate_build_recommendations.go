package part

import (
	"context"

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
	}
}

func (uc *GenerateBuildRecommendationsUseCase) Execute(ctx context.Context, args dto.GenerateBuildRecommendationsDTO) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	budgetCents := util.ConvertToCents(args.Budget)
	strategy := entity.GetStrategy(args.UsageType, budgetCents)
	allocations := strategy.GetAllocations()

	gpu, err := uc.selectBestGPUUC.Execute(ctx, SelectBestGPUArgs{
		BrandPreference: args.GpuPreference,
		MaxPriceCents:   int64(float64(budgetCents) * allocations[entity.TypeGPU]),
	})
	if err != nil {
		return nil, err
	}

	cpu, err := uc.selectBestCPUUC.Execute(ctx, SelectBestCPUArgs{
		BrandPreference: args.CpuPreference,
		MaxPriceCents:   int64(float64(budgetCents) * allocations[entity.TypeCPU]),
	})
	if err != nil {
		return nil, err
	}

	psu, err := uc.selectBestPSUUC.Execute(ctx, SelectBestPSUArgs{
		MinPSUWatts:   gpu.Specs.MinPSUWatts,
		MaxPriceCents: int64(float64(budgetCents) * allocations[entity.TypePSU]),
	})
	if err != nil {
		return nil, err
	}

	ram, err := uc.selectBestRAMUC.Execute(ctx, SelectBestRAMArgs{
		MaxPriceCents: int64(float64(budgetCents) * allocations[entity.TypeRAM]),
	})
	if err != nil {
		return nil, err
	}

	mobo, err := uc.selectBestMOBOUC.Execute(ctx, SelectBestMOBOArgs{
		Brand:         cpu.Brand,
		Socket:        cpu.Specs.Socket,
		MaxPriceCents: int64(float64(budgetCents) * allocations[entity.TypeMobo]),
	})
	if err != nil {
		return nil, err
	}

	var recommendedBuilds []presenters.RecommendationBuild
	recommendationBuild := presenters.RecommendationBuild{
		BuildType:  strategy.GetName(),
		Budget:     args.Budget,
		BuildValue: util.ConvertCentsToReal(gpu.PriceCents + cpu.PriceCents + psu.PriceCents + ram.PriceCents + mobo.PriceCents),
		Summary:    "High performance build suitable for gaming and productivity.",
		Parts: presenters.BuildParts{
			CPU:            &cpu,
			Motherboard:    &mobo,
			RAM:            &ram,
			GPU:            &gpu,
			PrimaryStorage: &entity.Part{},
			PSU:            &psu,
		},
	}

	recommendedBuilds = append(recommendedBuilds, recommendationBuild)

	return &presenters.RecommendedBuildsPresenter{
		Builds: recommendedBuilds,
	}, nil
}
