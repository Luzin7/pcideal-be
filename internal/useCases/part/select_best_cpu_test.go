package part

import (
	"context"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func newMockUpdatePartsUseCase() *UpdatePartsUseCase {
	return &UpdatePartsUseCase{}
}

func TestSelectBestCPU_SingleCPU(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	cpu := &entity.Part{
		ID:    primitive.NewObjectID(),
		Type:  entity.TypeCPU,
		Brand: "Intel",
		Model: "Core i7-12700K",
		Specs: entity.Specs{
			PerformanceScore: 85,
		},
		PriceCents: 250000,
		UpdatedAt:  time.Now(),
	}

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: []*entity.Part{cpu},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != cpu.Model {
		t.Errorf("Expected CPU model %s, got %s", cpu.Model, result.Model)
	}
}

func TestSelectBestCPU_SelectByPerformance(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i5-12600K",
			Specs: entity.Specs{
				PerformanceScore: 75,
			},
			PriceCents: 180000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i9-12900K",
			Specs: entity.Specs{
				PerformanceScore: 95,
			},
			PriceCents: 350000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: cpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Core i9-12900K" {
		t.Errorf("Expected best CPU to be Core i9-12900K, got %s", result.Model)
	}

	if result.Specs.PerformanceScore != 95 {
		t.Errorf("Expected performance score 95, got %d", result.Specs.PerformanceScore)
	}
}

func TestSelectBestCPU_SelectByPriceWhenSamePerformance(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 7 5800X",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 220000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 7 5800X3D",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 280000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: cpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Ryzen 7 5800X" {
		t.Errorf("Expected best CPU to be Ryzen 7 5800X (cheapest with same performance), got %s", result.Model)
	}

	if result.PriceCents != 220000 {
		t.Errorf("Expected price 220000, got %d", result.PriceCents)
	}
}

func TestSelectBestCPU_ComplexScenario(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i5-11400F",
			Specs: entity.Specs{
				PerformanceScore: 70,
			},
			PriceCents: 150000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 5 5600X",
			Specs: entity.Specs{
				PerformanceScore: 80,
			},
			PriceCents: 180000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 7 5800X",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 220000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: cpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Ryzen 7 5800X" {
		t.Errorf("Expected best CPU to be Ryzen 7 5800X, got %s", result.Model)
	}

	if result.Specs.PerformanceScore != 85 {
		t.Errorf("Expected performance score 85, got %d", result.Specs.PerformanceScore)
	}

	if result.PriceCents != 220000 {
		t.Errorf("Expected price 220000, got %d", result.PriceCents)
	}
}

func TestSelectBestCPU_EmptyList(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: []*entity.Part{},
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "" {
		t.Errorf("Expected empty CPU model, got %s", result.Model)
	}
}

func TestSelectBestCPU_AllSameSpecs(t *testing.T) {
	uc := NewSelectBestCPUUseCase(newMockUpdatePartsUseCase())

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K Store A",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K Store B",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 245000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i7-12700K Store C",
			Specs: entity.Specs{
				PerformanceScore: 85,
			},
			PriceCents: 240000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestCPUArgs{
		cpus: cpus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PriceCents != 240000 {
		t.Errorf("Expected cheapest price 240000, got %d", result.PriceCents)
	}

	if result.Model != "Core i7-12700K Store C" {
		t.Errorf("Expected Store C, got %s", result.Model)
	}
}
