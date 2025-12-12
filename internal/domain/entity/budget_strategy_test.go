package entity

import (
	"testing"
)

func TestGetStrategy_Gaming(t *testing.T) {
	tests := []struct {
		name         string
		budgetCents  int64
		expectedType string
		expectedName string
	}{
		{
			name:         "gaming low budget (below 3500 reais)",
			budgetCents:  300000,
			expectedType: "GamingLowBudgetStrategy",
			expectedName: "GAMING",
		},
		{
			name:         "gaming mid budget (between 3500 and 7000 reais)",
			budgetCents:  500000,
			expectedType: "GamingMidBudgetStrategy",
			expectedName: "GAMING",
		},
		{
			name:         "gaming high budget (above 7000 reais)",
			budgetCents:  800000,
			expectedType: "GamingHighBudgetStrategy",
			expectedName: "GAMING",
		},
		{
			name:         "gaming at low budget threshold",
			budgetCents:  350000,
			expectedType: "GamingMidBudgetStrategy",
			expectedName: "GAMING",
		},
		{
			name:         "gaming at high budget threshold",
			budgetCents:  700000,
			expectedType: "GamingHighBudgetStrategy",
			expectedName: "GAMING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetStrategy("GAMING", tt.budgetCents)
			if strategy.GetName() != tt.expectedName {
				t.Errorf("GetStrategy() name = %v, expected %v", strategy.GetName(), tt.expectedName)
			}
		})
	}
}

func TestGetStrategy_Work(t *testing.T) {
	tests := []struct {
		name         string
		budgetCents  int64
		expectedType string
		expectedName string
	}{
		{
			name:         "work low budget (below 3500 reais)",
			budgetCents:  300000,
			expectedType: "WorkLowBudgetStrategy",
			expectedName: "WORK",
		},
		{
			name:         "work mid budget (between 3500 and 7000 reais)",
			budgetCents:  500000,
			expectedType: "WorkMidBudgetStrategy",
			expectedName: "WORK",
		},
		{
			name:         "work high budget (above 7000 reais)",
			budgetCents:  800000,
			expectedType: "WorkHighBudgetStrategy",
			expectedName: "WORK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetStrategy("WORK", tt.budgetCents)
			if strategy.GetName() != tt.expectedName {
				t.Errorf("GetStrategy() name = %v, expected %v", strategy.GetName(), tt.expectedName)
			}
		})
	}
}

func TestGetStrategy_Office(t *testing.T) {
	tests := []struct {
		name         string
		budgetCents  int64
		expectedName string
	}{
		{
			name:         "office low budget",
			budgetCents:  200000,
			expectedName: "OFFICE",
		},
		{
			name:         "office high budget",
			budgetCents:  800000,
			expectedName: "OFFICE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetStrategy("OFFICE", tt.budgetCents)
			if strategy.GetName() != tt.expectedName {
				t.Errorf("GetStrategy() name = %v, expected %v", strategy.GetName(), tt.expectedName)
			}
		})
	}
}

func TestGamingLowBudgetStrategy_Allocations(t *testing.T) {
	strategy := GamingLowBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.40,
		TypeCPU:  0.20,
		TypeMobo: 0.14,
		TypeRAM:  0.10,
		TypePSU:  0.10,
		TypeSSD:  0.06,
	}

	verifyAllocations(t, allocations, expected, "GamingLowBudgetStrategy")
}

func TestGamingMidBudgetStrategy_Allocations(t *testing.T) {
	strategy := GamingMidBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.45,
		TypeCPU:  0.18,
		TypeMobo: 0.12,
		TypeRAM:  0.10,
		TypePSU:  0.08,
		TypeSSD:  0.07,
	}

	verifyAllocations(t, allocations, expected, "GamingMidBudgetStrategy")
}

func TestGamingHighBudgetStrategy_Allocations(t *testing.T) {
	strategy := GamingHighBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.50,
		TypeCPU:  0.15,
		TypeMobo: 0.12,
		TypeRAM:  0.10,
		TypePSU:  0.07,
		TypeSSD:  0.06,
	}

	verifyAllocations(t, allocations, expected, "GamingHighBudgetStrategy")
}

func TestWorkLowBudgetStrategy_Allocations(t *testing.T) {
	strategy := WorkLowBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.45,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.15,
		TypePSU:  0.10,
	}

	verifyAllocations(t, allocations, expected, "WorkLowBudgetStrategy")
}

func TestWorkMidBudgetStrategy_Allocations(t *testing.T) {
	strategy := WorkMidBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.40,
		TypeRAM:  0.25,
		TypeSSD:  0.12,
		TypeMobo: 0.15,
		TypePSU:  0.08,
	}

	verifyAllocations(t, allocations, expected, "WorkMidBudgetStrategy")
}

func TestWorkHighBudgetStrategy_Allocations(t *testing.T) {
	strategy := WorkHighBudgetStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.25,
		TypeCPU:  0.35,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.05,
		TypePSU:  0.05,
	}

	verifyAllocations(t, allocations, expected, "WorkHighBudgetStrategy")
}

func TestOfficeStrategy_Allocations(t *testing.T) {
	strategy := OfficeStrategy{}
	allocations := strategy.GetAllocations()

	expected := map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.40,
		TypeRAM:  0.15,
		TypeSSD:  0.15,
		TypeMobo: 0.15,
		TypePSU:  0.15,
	}

	verifyAllocations(t, allocations, expected, "OfficeStrategy")
}

func TestAllocations_SumToOne(t *testing.T) {
	strategies := []BudgetStrategy{
		GamingLowBudgetStrategy{},
		GamingMidBudgetStrategy{},
		GamingHighBudgetStrategy{},
		WorkLowBudgetStrategy{},
		WorkMidBudgetStrategy{},
		WorkHighBudgetStrategy{},
		OfficeStrategy{},
	}

	for _, strategy := range strategies {
		t.Run(strategy.GetName()+" allocations sum", func(t *testing.T) {
			allocations := strategy.GetAllocations()
			sum := 0.0
			for _, value := range allocations {
				sum += value
			}

			if sum < 0.99 || sum > 1.01 {
				t.Errorf("%s allocations sum to %f, expected approximately 1.0", strategy.GetName(), sum)
			}
		})
	}
}

func verifyAllocations(t *testing.T, actual, expected map[PartType]float64, strategyName string) {
	for partType, expectedValue := range expected {
		actualValue, exists := actual[partType]
		if !exists {
			t.Errorf("%s missing allocation for %s", strategyName, partType)
			continue
		}
		if actualValue != expectedValue {
			t.Errorf("%s allocation for %s = %f, expected %f", strategyName, partType, actualValue, expectedValue)
		}
	}

	for partType := range actual {
		if _, exists := expected[partType]; !exists {
			t.Errorf("%s has unexpected allocation for %s", strategyName, partType)
		}
	}
}

func TestGetStrategy_DefaultsToGaming(t *testing.T) {
	tests := []struct {
		name        string
		objective   string
		budgetCents int64
	}{
		{
			name:        "empty string defaults to gaming",
			objective:   "",
			budgetCents: 500000,
		},
		{
			name:        "unknown objective defaults to gaming",
			objective:   "UNKNOWN",
			budgetCents: 500000,
		},
		{
			name:        "random string defaults to gaming",
			objective:   "RANDOM",
			budgetCents: 500000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy := GetStrategy(tt.objective, tt.budgetCents)
			if strategy.GetName() != "GAMING" {
				t.Errorf("GetStrategy(%q) name = %v, expected GAMING", tt.objective, strategy.GetName())
			}
		})
	}
}
