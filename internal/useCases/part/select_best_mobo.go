package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestMOBOUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestMOBOUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestMOBOUseCase {
	return &SelectBestMOBOUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestMOBOArgs struct {
	mobos []*entity.Part
}

func (uc *SelectBestMOBOUseCase) Execute(ctx context.Context, args SelectBestMOBOArgs) (entity.Part, *errors.ErrService) {
	queueStalePartsUpdate(uc.UpdatePartsUseCase, args.mobos)

	bestMOBO := chooseBestPart(args.mobos, partRanking{
		score: func(part *entity.Part) int { return int(part.Specs.TierScore) },
	})
	if bestMOBO == nil {
		return entity.Part{}, nil
	}

	return *bestMOBO, nil
}
