package entity

type BudgetStrategy interface {
	GetAllocations() map[PartType]float64
	GetName() string
}

type GamingLowBudgetStrategy struct{}

func (s GamingLowBudgetStrategy) GetName() string { return "GAMING" }
func (s GamingLowBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.35,
		TypeCPU:  0.25,
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
		TypeCPU:  0.30,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.10,
		TypePSU:  0.05,
	}
}

type ContentCreatorLowBudgetStrategy struct{}

func (s ContentCreatorLowBudgetStrategy) GetName() string { return "CONTENT_CREATOR" }
func (s ContentCreatorLowBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeCPU:  0.35,
		TypeGPU:  0.20,
		TypeRAM:  0.15,
		TypeMobo: 0.14,
		TypeSSD:  0.08,
		TypePSU:  0.08,
	}
}

type ContentCreatorMidBudgetStrategy struct{}

func (s ContentCreatorMidBudgetStrategy) GetName() string { return "CONTENT_CREATOR" }
func (s ContentCreatorMidBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.30,
		TypeCPU:  0.28,
		TypeRAM:  0.15,
		TypeMobo: 0.12,
		TypeSSD:  0.08,
		TypePSU:  0.07,
	}
}

type ContentCreatorHighBudgetStrategy struct{}

func (s ContentCreatorHighBudgetStrategy) GetName() string { return "CONTENT_CREATOR" }
func (s ContentCreatorHighBudgetStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.35,
		TypeCPU:  0.25,
		TypeRAM:  0.15,
		TypeMobo: 0.12,
		TypeSSD:  0.08,
		TypePSU:  0.05,
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

	case "CONTENT_CREATOR":
		if budget < lowBudgetThreshold {
			return ContentCreatorLowBudgetStrategy{}
		} else if budget < highBudgetThreshold {
			return ContentCreatorMidBudgetStrategy{}
		}
		return ContentCreatorHighBudgetStrategy{}

	default: // "GAMING"
		if budget < lowBudgetThreshold {
			return GamingLowBudgetStrategy{}
		} else if budget < highBudgetThreshold {
			return GamingMidBudgetStrategy{}
		}
		return GamingHighBudgetStrategy{}
	}
}
