package part

import (
	"context"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIntegration_CompleteComponentSelection(t *testing.T) {
	updateUC := newMockUpdatePartsUseCase()
	cpuUC := NewSelectBestCPUUseCase(updateUC)
	gpuUC := NewSelectBestGPUUseCase(updateUC)
	ramUC := NewSelectBestRAMUseCase(updateUC)
	psuUC := NewSelectBestPSUUseCase(updateUC)
	moboUC := NewSelectBestMOBOUseCase(updateUC)

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 5 5600X",
			Specs: entity.Specs{
				Socket:           "AM4",
				PerformanceScore: 80,
			},
			PriceCents: 180000,
			UpdatedAt:  time.Now(),
		},
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "Intel",
			Model: "Core i5-12400F",
			Specs: entity.Specs{
				Socket:           "LGA1700",
				PerformanceScore: 75,
			},
			PriceCents: 170000,
			UpdatedAt:  time.Now(),
		},
	}

	gpus := []*entity.Part{
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
				PerformanceScore: 82,
				VramGB:           12,
			},
			PriceCents: 310000,
			UpdatedAt:  time.Now(),
		},
	}

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
	}

	psus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypePSU,
			Brand: "Corsair",
			Model: "CV650 650W 80+ Bronze",
			Specs: entity.Specs{
				Wattage:          650,
				EfficiencyRating: 80,
			},
			PriceCents: 40000,
			UpdatedAt:  time.Now(),
		},
	}

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
	}

	selectedCPU, errCPU := cpuUC.Execute(context.Background(), SelectBestCPUArgs{cpus: cpus})
	selectedGPU, errGPU := gpuUC.Execute(context.Background(), SelectBestGPUArgs{gpus: gpus})
	selectedRAM, errRAM := ramUC.Execute(context.Background(), SelectBestRAMArgs{rams: rams})
	selectedPSU, errPSU := psuUC.Execute(context.Background(), SelectBestPSUArgs{psus: psus})
	selectedMOBO, errMOBO := moboUC.Execute(context.Background(), SelectBestMOBOArgs{mobos: mobos})

	if errCPU != nil {
		t.Fatalf("Error selecting CPU: %v", errCPU)
	}
	if errGPU != nil {
		t.Fatalf("Error selecting GPU: %v", errGPU)
	}
	if errRAM != nil {
		t.Fatalf("Error selecting RAM: %v", errRAM)
	}
	if errPSU != nil {
		t.Fatalf("Error selecting PSU: %v", errPSU)
	}
	if errMOBO != nil {
		t.Fatalf("Error selecting MOBO: %v", errMOBO)
	}

	t.Run("CPU Selection", func(t *testing.T) {
		if selectedCPU.Model != "Ryzen 5 5600X" {
			t.Errorf("Expected Ryzen 5 5600X (better performance), got %s", selectedCPU.Model)
		}
	})

	t.Run("GPU Selection", func(t *testing.T) {
		if selectedGPU.Model != "RX 6700 XT" {
			t.Errorf("Expected RX 6700 XT (better performance), got %s", selectedGPU.Model)
		}
	})

	t.Run("RAM Selection", func(t *testing.T) {
		if selectedRAM.Specs.CasLatency != 16 {
			t.Errorf("Expected CAS Latency 16, got %d", selectedRAM.Specs.CasLatency)
		}
	})

	t.Run("PSU Selection", func(t *testing.T) {
		if selectedPSU.Specs.Wattage != 650 {
			t.Errorf("Expected 650W PSU, got %dW", selectedPSU.Specs.Wattage)
		}
	})

	t.Run("MOBO Selection", func(t *testing.T) {
		if selectedMOBO.Specs.Socket != "AM4" {
			t.Errorf("Expected AM4 socket, got %s", selectedMOBO.Specs.Socket)
		}
	})

	t.Run("Total Build Cost", func(t *testing.T) {
		totalCents := selectedCPU.PriceCents +
			selectedGPU.PriceCents +
			selectedRAM.PriceCents +
			selectedPSU.PriceCents +
			selectedMOBO.PriceCents

		expectedTotal := int64(600000)
		if totalCents < expectedTotal-100000 || totalCents > expectedTotal+100000 {
			t.Errorf("Total build cost %d is outside expected range around %d", totalCents, expectedTotal)
		}

		t.Logf("Build Summary:")
		t.Logf("  CPU:  %s - R$ %.2f", selectedCPU.Model, float64(selectedCPU.PriceCents)/100)
		t.Logf("  GPU:  %s - R$ %.2f", selectedGPU.Model, float64(selectedGPU.PriceCents)/100)
		t.Logf("  RAM:  %s - R$ %.2f", selectedRAM.Model, float64(selectedRAM.PriceCents)/100)
		t.Logf("  PSU:  %s - R$ %.2f", selectedPSU.Model, float64(selectedPSU.PriceCents)/100)
		t.Logf("  MOBO: %s - R$ %.2f", selectedMOBO.Model, float64(selectedMOBO.PriceCents)/100)
		t.Logf("  TOTAL: R$ %.2f", float64(totalCents)/100)
	})
}

func TestIntegration_CompatibilityCheck(t *testing.T) {
	updateUC := newMockUpdatePartsUseCase()
	cpuUC := NewSelectBestCPUUseCase(updateUC)
	moboUC := NewSelectBestMOBOUseCase(updateUC)

	cpus := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeCPU,
			Brand: "AMD",
			Model: "Ryzen 7 5800X",
			Specs: entity.Specs{
				Socket:           "AM4",
				PerformanceScore: 85,
			},
			PriceCents: 250000,
			UpdatedAt:  time.Now(),
		},
	}

	mobosAM4 := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "ASUS",
			Model: "B550-Plus AM4",
			Specs: entity.Specs{
				Socket:    "AM4",
				TierScore: 7,
			},
			PriceCents: 80000,
			UpdatedAt:  time.Now(),
		},
	}

	mobosLGA1700 := []*entity.Part{
		{
			ID:    primitive.NewObjectID(),
			Type:  entity.TypeMobo,
			Brand: "ASUS",
			Model: "B660-Plus LGA1700",
			Specs: entity.Specs{
				Socket:    "LGA1700",
				TierScore: 7,
			},
			PriceCents: 85000,
			UpdatedAt:  time.Now(),
		},
	}

	selectedCPU, _ := cpuUC.Execute(context.Background(), SelectBestCPUArgs{cpus: cpus})

	selectedMOBOCompatible, _ := moboUC.Execute(context.Background(), SelectBestMOBOArgs{mobos: mobosAM4})

	selectedMOBOIncompatible, _ := moboUC.Execute(context.Background(), SelectBestMOBOArgs{mobos: mobosLGA1700})

	t.Run("Compatible Socket", func(t *testing.T) {
		if selectedCPU.Specs.Socket != selectedMOBOCompatible.Specs.Socket {
			t.Errorf("CPU socket %s should match MOBO socket %s",
				selectedCPU.Specs.Socket, selectedMOBOCompatible.Specs.Socket)
		}
	})

	t.Run("Incompatible Socket", func(t *testing.T) {
		if selectedCPU.Specs.Socket == selectedMOBOIncompatible.Specs.Socket {
			t.Errorf("Expected incompatible sockets, but both are %s", selectedCPU.Specs.Socket)
		}
		t.Logf("Correctly identified incompatible sockets: CPU=%s, MOBO=%s",
			selectedCPU.Specs.Socket, selectedMOBOIncompatible.Specs.Socket)
	})
}

func TestIntegration_BudgetAllocation(t *testing.T) {
	totalBudgetCents := int64(500000)

	strategy := entity.GamingMidBudgetStrategy{}
	allocations := strategy.GetAllocations()

	budgetBreakdown := make(map[entity.PartType]int64)
	for partType, percentage := range allocations {
		budgetBreakdown[partType] = int64(float64(totalBudgetCents) * percentage)
	}

	t.Run("Budget Allocations", func(t *testing.T) {
		t.Logf("Total Budget: R$ %.2f", float64(totalBudgetCents)/100)
		t.Logf("Strategy: %s", strategy.GetName())
		t.Logf("")
		t.Logf("Component Budgets:")

		for partType, budget := range budgetBreakdown {
			percentage := allocations[partType] * 100
			t.Logf("  %s: R$ %.2f (%.1f%%)", partType, float64(budget)/100, percentage)
		}

		if budgetBreakdown[entity.TypeGPU] <= budgetBreakdown[entity.TypeCPU] {
			t.Error("Gaming strategy should allocate more budget to GPU than CPU")
		}

		sum := int64(0)
		for _, budget := range budgetBreakdown {
			sum += budget
		}

		tolerance := int64(100)
		if sum < totalBudgetCents-tolerance || sum > totalBudgetCents+tolerance {
			t.Errorf("Budget sum %d should be close to total %d", sum, totalBudgetCents)
		}
	})
}
