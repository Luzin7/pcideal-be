package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestGPUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestGPUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestGPUUseCase {
	return &SelectBestGPUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestGPUArgs struct {
	gpus []*entity.Part
}

func (uc *SelectBestGPUUseCase) Execute(ctx context.Context, args SelectBestGPUArgs) (entity.Part, *errors.ErrService) {
	queueStalePartsUpdate(uc.UpdatePartsUseCase, args.gpus)

	bestGPU := chooseBestPart(args.gpus, partRanking{score: performanceScore})
	if bestGPU == nil {
		return entity.Part{}, nil
	}

	return *bestGPU, nil
}
