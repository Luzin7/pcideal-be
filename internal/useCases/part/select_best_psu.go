package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestPSUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestPSUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestPSUUseCase {
	return &SelectBestPSUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestPSUArgs struct {
	psus []*entity.Part
}

func (uc *SelectBestPSUUseCase) Execute(ctx context.Context, args SelectBestPSUArgs) (entity.Part, *errors.ErrService) {
	queueStalePartsUpdate(uc.UpdatePartsUseCase, args.psus)

	bestPSU := chooseBestPart(args.psus, partRanking{
		score: func(part *entity.Part) int { return int(part.Specs.EfficiencyRating) },
	})
	if bestPSU == nil {
		return entity.Part{}, nil
	}

	return *bestPSU, nil
}
