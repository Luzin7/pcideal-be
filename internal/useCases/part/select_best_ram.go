package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestRAMUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestRAMUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestRAMUseCase {
	return &SelectBestRAMUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestRAMArgs struct {
	rams []*entity.Part
}

func (uc *SelectBestRAMUseCase) Execute(ctx context.Context, args SelectBestRAMArgs) (entity.Part, *errors.ErrService) {
	validRAMs := make([]*entity.Part, 0, len(args.rams))

	for _, ram := range args.rams {
		if ram.Specs.CasLatency <= 0 {
			continue
		}
		validRAMs = append(validRAMs, ram)
	}

	queueStalePartsUpdate(uc.UpdatePartsUseCase, validRAMs)

	bestRAM := chooseBestPart(validRAMs, partRanking{
		score:    performanceScore,
		lowerTie: func(part *entity.Part) int { return part.Specs.CasLatency },
	})
	if bestRAM == nil {
		return entity.Part{}, nil
	}

	return *bestRAM, nil
}
