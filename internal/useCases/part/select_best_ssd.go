package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestSSDUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestSSDUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestSSDUseCase {
	return &SelectBestSSDUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestSSDArgs struct {
	ssds []*entity.Part
}

func (uc *SelectBestSSDUseCase) Execute(ctx context.Context, args SelectBestSSDArgs) (entity.Part, *errors.ErrService) {
	queueStalePartsUpdate(uc.UpdatePartsUseCase, args.ssds)

	bestSSD := chooseBestPart(args.ssds, partRanking{score: performanceScore})
	if bestSSD == nil {
		return entity.Part{}, nil
	}

	return *bestSSD, nil
}
