package part

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/dto"
)

type emptyPartRepository struct{}

func (emptyPartRepository) GetAllParts(ctx context.Context, args repository.FindAllPartsArgs) ([]*entity.Part, error) {
	return nil, nil
}

func (emptyPartRepository) GetPartByID(ctx context.Context, id string) (*entity.Part, error) {
	return nil, nil
}

func (emptyPartRepository) UpdatePart(ctx context.Context, partId string, part *entity.Part) error {
	return nil
}

func (emptyPartRepository) FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args repository.FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error) {
	return nil, nil
}

type catalogPartRepository struct {
	emptyPartRepository
	parts []*entity.Part
}

func (r catalogPartRepository) FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args repository.FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error) {
	var matches []*entity.Part

	for _, part := range r.parts {
		if part.Type != args.PartType || part.PriceCents > args.MaxPriceCents {
			continue
		}
		matches = append(matches, part)
	}

	return matches, nil
}

func buildCatalog() []*entity.Part {
	now := time.Now()

	return []*entity.Part{
		{
			Type:       entity.TypeCPU,
			Brand:      "AMD",
			Model:      "Ryzen 5 5600",
			PriceCents: 85000,
			Specs:      entity.Specs{Socket: "AM4", MemoryType: "DDR4", PerformanceScore: 80, MinPSUWatts: 450},
			UpdatedAt:  now,
		},
		{
			Type:       entity.TypeMobo,
			Brand:      "ASUS",
			Model:      "B550-Plus",
			PriceCents: 50000,
			Specs:      entity.Specs{Socket: "AM4", MemoryType: "DDR4", TierScore: 7},
			UpdatedAt:  now,
		},
		{
			Type:       entity.TypeGPU,
			Brand:      "NVIDIA",
			Model:      "RTX 4060",
			PriceCents: 200000,
			Specs:      entity.Specs{PerformanceScore: 80, VramGB: 8},
			UpdatedAt:  now,
		},
		{
			Type:       entity.TypePSU,
			Brand:      "Corsair",
			Model:      "CV650",
			PriceCents: 30000,
			Specs:      entity.Specs{Wattage: 650, EfficiencyRating: 80},
			UpdatedAt:  now,
		},
		{
			Type:       entity.TypeRAM,
			Brand:      "Corsair",
			Model:      "Vengeance 16GB",
			PriceCents: 40000,
			Specs:      entity.Specs{MemoryType: "DDR4", PerformanceScore: 75, CasLatency: 16},
			UpdatedAt:  now,
		},
		{
			Type:       entity.TypeSSD,
			Brand:      "Kingston",
			Model:      "NV2 1TB",
			PriceCents: 30000,
			Specs:      entity.Specs{CapacityGB: 1000},
			UpdatedAt:  now,
		},
	}
}

func TestGenerateBuildRecommendations_SkipsUnbuildableTier(t *testing.T) {
	updatePartsUC := &UpdatePartsUseCase{}

	uc := NewGenerateBuildRecommendationsUseCase(
		catalogPartRepository{parts: buildCatalog()},
		updatePartsUC,
		NewSelectBestCPUUseCase(updatePartsUC),
		NewSelectBestGPUUseCase(updatePartsUC),
		NewSelectBestPSUUseCase(updatePartsUC),
		NewSelectBestRAMUseCase(updatePartsUC),
		NewSelectBestMOBOUseCase(updatePartsUC),
		NewSelectBestSSDUseCase(updatePartsUC),
	)

	result, err := uc.Execute(context.Background(), dto.GenerateBuildRecommendationsDTO{
		UsageType:   "GAMING",
		BudgetCents: 500000,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result.Builds) != 2 {
		t.Fatalf("expected 2 builds, got %d", len(result.Builds))
	}
	if result.Builds[0].BuildType != "BALANCED" {
		t.Fatalf("expected first build BALANCED, got %s", result.Builds[0].BuildType)
	}
	if result.Builds[0].Parts.CPU.Model != "Ryzen 5 5600" {
		t.Fatalf("expected BALANCED CPU, got %q", result.Builds[0].Parts.CPU.Model)
	}
}

func TestGenerateBuildRecommendations_EmptyCandidates(t *testing.T) {
	updatePartsUC := &UpdatePartsUseCase{}

	uc := NewGenerateBuildRecommendationsUseCase(
		emptyPartRepository{},
		updatePartsUC,
		NewSelectBestCPUUseCase(updatePartsUC),
		NewSelectBestGPUUseCase(updatePartsUC),
		NewSelectBestPSUUseCase(updatePartsUC),
		NewSelectBestRAMUseCase(updatePartsUC),
		NewSelectBestMOBOUseCase(updatePartsUC),
		NewSelectBestSSDUseCase(updatePartsUC),
	)

	result, err := uc.Execute(context.Background(), dto.GenerateBuildRecommendationsDTO{
		UsageType:   "GAMING",
		BudgetCents: 500000,
	})

	if result != nil {
		t.Fatalf("expected no build result, got %+v", result)
	}
	if err == nil {
		t.Fatal("expected no-compatible-parts error, got nil")
	}
	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status %d, got %d", http.StatusUnprocessableEntity, err.StatusCode)
	}
	if !strings.Contains(err.Message, "CPU") {
		t.Errorf("expected message to name the missing part type, got %q", err.Message)
	}
}
