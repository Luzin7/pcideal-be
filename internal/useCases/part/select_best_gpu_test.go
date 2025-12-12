package part

import (
	"context"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSelectBestGPU_SingleGPU(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	gpu := &entity.Part{
		ID:    primitive.NewObjectID(),
		Type:  entity.TypeGPU,
		Brand: "NVIDIA",
		Model: "RTX 3070",
		Specs: entity.Specs{
			PerformanceScore: 85,
			VramGB:           8,
		},
		PriceCents: 350000,
		UpdatedAt:  time.Now(),
	}

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: []*entity.Part{gpu},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != gpu.Model {
		t.Errorf("Expected GPU model %s, got %s", gpu.Model, result.Model)
	}
}

func TestSelectBestGPU_SelectByPerformance(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	gpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3060",
			Specs: entity.Specs{
				PerformanceScore: 75,
				VramGB:           12,
			},
			PriceCents: 280000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 350000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3080",
			Specs: entity.Specs{
				PerformanceScore: 95,
				VramGB:           10,
			},
			PriceCents: 450000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: gpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "RTX 3080" {
		t.Errorf("Expected best GPU to be RTX 3080, got %s", result.Model)
	}

	if result.Specs.PerformanceScore != 95 {
		t.Errorf("Expected performance score 95, got %d", result.Specs.PerformanceScore)
	}
}

func TestSelectBestGPU_SelectByPriceWhenSamePerformance(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	gpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 350000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "AMD",
			Model: "RX 6800",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           16,
			},
			PriceCents: 320000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070 Ti",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 380000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: gpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "RX 6800" {
		t.Errorf("Expected best GPU to be RX 6800 (cheapest with same performance), got %s", result.Model)
	}

	if result.PriceCents != 320000 {
		t.Errorf("Expected price 320000, got %d", result.PriceCents)
	}
}

func TestSelectBestGPU_ComplexScenario(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	gpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "GTX 1660 Super",
			Specs: entity.Specs{
				PerformanceScore: 65,
				VramGB:           6,
			},
			PriceCents: 200000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "AMD",
			Model: "RX 6600 XT",
			Specs: entity.Specs{
				PerformanceScore: 75,
				VramGB:           8,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3060 Ti",
			Specs: entity.Specs{
				PerformanceScore: 80,
				VramGB:           8,
			},
			PriceCents: 300000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "AMD",
			Model: "RX 6700 XT",
			Specs: entity.Specs{
				PerformanceScore: 80,
				VramGB:           12,
			},
			PriceCents: 280000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: gpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "RX 6700 XT" {
		t.Errorf("Expected best GPU to be RX 6700 XT, got %s", result.Model)
	}

	if result.Specs.PerformanceScore != 80 {
		t.Errorf("Expected performance score 80, got %d", result.Specs.PerformanceScore)
	}

	if result.PriceCents != 280000 {
		t.Errorf("Expected price 280000, got %d", result.PriceCents)
	}
}

func TestSelectBestGPU_EmptyList(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: []*entity.Part{},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "" {
		t.Errorf("Expected empty GPU model, got %s", result.Model)
	}
}

func TestSelectBestGPU_AllSameSpecs(t *testing.T) {
	uc := NewSelectBestGPUUseCase(newMockUpdatePartsUseCase())

	gpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070 Store A",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 360000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070 Store B",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 345000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeGPU,
			Brand: "NVIDIA",
			Model: "RTX 3070 Store C",
			Specs: entity.Specs{
				PerformanceScore: 85,
				VramGB:           8,
			},
			PriceCents: 350000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestGPUArgs{
		gpus: gpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PriceCents != 345000 {
		t.Errorf("Expected cheapest price 345000, got %d", result.PriceCents)
	}

	if result.Model != "RTX 3070 Store B" {
		t.Errorf("Expected Store B, got %s", result.Model)
	}
}
