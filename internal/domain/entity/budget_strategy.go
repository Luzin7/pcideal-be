package entity

type BudgetStrategy interface {
	GetAllocations() map[PartType]float64
	GetName() string
}

type GamingLowBudgetStrategy struct{}

func (s GamingLowBudgetStrategy) GetName() string { return "GAMING" }
func (s GamingLowBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.40,
		TypeCPU:  0.20,
		TypeMobo: 0.14,
		TypeRAM:  0.10,
		TypePSU:  0.10,
		TypeSSD:  0.06,
	}
}

type GamingMidBudgetStrategy struct{}

func (s GamingMidBudgetStrategy) GetName() string { return "GAMING" }
func (s GamingMidBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.45,
		TypeCPU:  0.18,
		TypeMobo: 0.12,
		TypeRAM:  0.10,
		TypePSU:  0.08,
		TypeSSD:  0.07,
	}
}

type GamingHighBudgetStrategy struct{}

func (s GamingHighBudgetStrategy) GetName() string { return "GAMING" }
func (s GamingHighBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.50,
		TypeCPU:  0.15,
		TypeMobo: 0.12,
		TypeRAM:  0.10,
		TypePSU:  0.07,
		TypeSSD:  0.06,
	}
}

type WorkLowBudgetStrategy struct{}

func (s WorkLowBudgetStrategy) GetName() string { return "WORK" }
func (s WorkLowBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.45,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.15,
		TypePSU:  0.10,
	}
}

type WorkMidBudgetStrategy struct{}

func (s WorkMidBudgetStrategy) GetName() string { return "WORK" }
func (s WorkMidBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.40,
		TypeRAM:  0.25,
		TypeSSD:  0.12,
		TypeMobo: 0.15,
		TypePSU:  0.08,
	}
}

type WorkHighBudgetStrategy struct{}

func (s WorkHighBudgetStrategy) GetName() string { return "WORK" }
func (s WorkHighBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.25,
		TypeCPU:  0.35,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.05,
		TypePSU:  0.05,
	}
}

type OfficeStrategy struct{}

func (s OfficeStrategy) GetName() string { return "OFFICE" }
func (s OfficeStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.00,
		TypeCPU:  0.40,
		TypeRAM:  0.15,
		TypeSSD:  0.15,
		TypeMobo: 0.15,
		TypePSU:  0.15,
	}
}

func GetStrategy(objective string, budget int64) BudgetStrategy {
	const (
		lowBudgetThreshold  = 350000
		highBudgetThreshold = 700000
	)

	switch objective {
	case "WORK":
		if budget < lowBudgetThreshold {
			return WorkLowBudgetStrategy{}
		} else if budget < highBudgetThreshold {
			return WorkMidBudgetStrategy{}
		}
		return WorkHighBudgetStrategy{}
	case "OFFICE":
		return OfficeStrategy{}
	default: // "GAMING"
		if budget < lowBudgetThreshold {
			return GamingLowBudgetStrategy{}
		} else if budget < highBudgetThreshold {
			return GamingMidBudgetStrategy{}
		}
		return GamingHighBudgetStrategy{}
	}
}
