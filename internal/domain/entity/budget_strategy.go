package entity

type BudgetStrategy interface {
	GetAllocations() map[PartType]float64
	GetName() string
}

type GamerStrategy struct{}

func (s GamerStrategy) GetName() string { return "GAMING" }
func (s GamerStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.45,
		TypeCPU:  0.20,
		TypeMobo: 0.10,
		TypeRAM:  0.08,
		TypePSU:  0.07,
		TypeSSD:  0.05,
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

type WorkstationStrategy struct{}

func (s WorkstationStrategy) GetName() string { return "WORK" }
func (s WorkstationStrategy) GetAllocations() map[PartType]float64 {
	return map[PartType]float64{
		TypeGPU:  0.20,
		TypeCPU:  0.35,
		TypeRAM:  0.20,
		TypeSSD:  0.10,
		TypeMobo: 0.10,
		TypePSU:  0.05,
	}
}

func GetStrategy(objective string) BudgetStrategy {
	switch objective {
	case "WORK":
		return WorkstationStrategy{}
	case "OFFICE":
		return OfficeStrategy{}
	default:
		return GamerStrategy{}
	}
}
