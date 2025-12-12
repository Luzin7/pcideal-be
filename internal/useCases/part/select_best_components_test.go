package part

import (
	"context"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSelectBestRAM_ByPerformanceScore(t *testing.T) {
	uc := NewSelectBestRAMUseCase(newMockUpdatePartsUseCase())

	rams := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Corsair",
			Model: "Vengeance 16GB 3200MHz CL16",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       16,
			},
			PriceCents: 30000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "G.Skill",
			Model: "Trident Z 16GB 3600MHz CL18",
			Specs: entity.Specs{
				PerformanceScore: 80,
				MemorySpeedMHz:   3600,
				CasLatency:       18,
			},
			PriceCents: 35000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestRAMArgs{
		rams: rams,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Trident Z 16GB 3600MHz CL18" {
		t.Errorf("Expected best RAM to be Trident Z (higher performance), got %s", result.Model)
	}
}

func TestSelectBestRAM_ByCasLatencyWhenSamePerformance(t *testing.T) {
	uc := NewSelectBestRAMUseCase(newMockUpdatePartsUseCase())

	rams := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Corsair",
			Model: "Vengeance 16GB 3200MHz CL18",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       18,
			},
			PriceCents: 30000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Kingston",
			Model: "Fury 16GB 3200MHz CL16",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       16,
			},
			PriceCents: 32000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestRAMArgs{
		rams: rams,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Fury 16GB 3200MHz CL16" {
		t.Errorf("Expected best RAM to be Fury (lower CAS latency), got %s", result.Model)
	}
}

func TestSelectBestRAM_ByPriceWhenAllSame(t *testing.T) {
	uc := NewSelectBestRAMUseCase(newMockUpdatePartsUseCase())

	rams := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Corsair",
			Model: "Vengeance Store A",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       16,
			},
			PriceCents: 35000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Corsair",
			Model: "Vengeance Store B",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       16,
			},
			PriceCents: 30000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestRAMArgs{
		rams: rams,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PriceCents != 30000 {
		t.Errorf("Expected cheapest price 30000, got %d", result.PriceCents)
	}
}

func TestSelectBestRAM_SkipInvalidCasLatency(t *testing.T) {
	uc := NewSelectBestRAMUseCase(newMockUpdatePartsUseCase())

	rams := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Corsair",
			Model: "Invalid RAM",
			Specs: entity.Specs{
				PerformanceScore: 90,
				MemorySpeedMHz:   3600,
				CasLatency:       0,
			},
			PriceCents: 20000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeRAM,
			Brand: "Kingston",
			Model: "Valid RAM",
			Specs: entity.Specs{
				PerformanceScore: 75,
				MemorySpeedMHz:   3200,
				CasLatency:       16,
			},
			PriceCents: 30000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestRAMArgs{
		rams: rams,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "Valid RAM" {
		t.Errorf("Expected Valid RAM (skipped invalid CAS latency), got %s", result.Model)
	}
}

func TestSelectBestPSU_ByEfficiencyRating(t *testing.T) {
	uc := NewSelectBestPSUUseCase(newMockUpdatePartsUseCase())

	psus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypePSU,
			Brand: "Corsair",
			Model: "CV650 Bronze",
			Specs: entity.Specs{
				Wattage:          650,
				EfficiencyRating: 80,
			},
			PriceCents: 40000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypePSU,
			Brand: "EVGA",
			Model: "650 G3 Gold",
			Specs: entity.Specs{
				Wattage:          650,
				EfficiencyRating: 90,
			},
			PriceCents: 50000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestPSUArgs{
		psus: psus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "650 G3 Gold" {
		t.Errorf("Expected best PSU to be 650 G3 Gold (higher efficiency), got %s", result.Model)
	}
}

func TestSelectBestPSU_ByPriceWhenSameEfficiency(t *testing.T) {
	uc := NewSelectBestPSUUseCase(newMockUpdatePartsUseCase())

	psus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypePSU,
			Brand: "Corsair",
			Model: "RM650x Gold Store A",
			Specs: entity.Specs{
				Wattage:          650,
				EfficiencyRating: 90,
			},
			PriceCents: 55000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypePSU,
			Brand: "Corsair",
			Model: "RM650x Gold Store B",
			Specs: entity.Specs{
				Wattage:          650,
				EfficiencyRating: 90,
			},
			PriceCents: 50000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestPSUArgs{
		psus: psus,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PriceCents != 50000 {
		t.Errorf("Expected cheapest price 50000, got %d", result.PriceCents)
	}
}

func TestSelectBestMOBO_ByTierScore(t *testing.T) {
	uc := NewSelectBestMOBOUseCase(newMockUpdatePartsUseCase())

	mobos := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "ASUS",
			Model: "B550-Plus",
			Specs: entity.Specs{
				Socket:    "AM4",
				TierScore: 7,
			},
			PriceCents: 80000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "MSI",
			Model: "X570 Tomahawk",
			Specs: entity.Specs{
				Socket:    "AM4",
				TierScore: 9,
			},
			PriceCents: 120000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestMOBOArgs{
		mobos: mobos,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Model != "X570 Tomahawk" {
		t.Errorf("Expected best MOBO to be X570 Tomahawk (higher tier), got %s", result.Model)
	}
}

func TestSelectBestMOBO_ByPriceWhenSameTier(t *testing.T) {
	uc := NewSelectBestMOBOUseCase(newMockUpdatePartsUseCase())

	mobos := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "ASUS",
			Model: "B550-Plus Store A",
			Specs: entity.Specs{
				Socket:    "AM4",
				TierScore: 7,
			},
			PriceCents: 85000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "ASUS",
			Model: "B550-Plus Store B",
			Specs: entity.Specs{
				Socket:    "AM4",
				TierScore: 7,
			},
			PriceCents: 80000,
			UpdatedAt:  time.Now(),
		},
	}

	result, err := uc.Execute(context.Background(), SelectBestMOBOArgs{
		mobos: mobos,
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.PriceCents != 80000 {
		t.Errorf("Expected cheapest price 80000, got %d", result.PriceCents)
	}
}
