package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type SelectBestCPUUseCase struct {
	UpdatePartsUseCase *UpdatePartsUseCase
}

func NewSelectBestCPUUseCase(updatePartsUseCase *UpdatePartsUseCase) *SelectBestCPUUseCase {
	return &SelectBestCPUUseCase{
		UpdatePartsUseCase: updatePartsUseCase,
	}
}

type SelectBestCPUArgs struct {
	cpus []*entity.Part
}

func (uc *SelectBestCPUUseCase) Execute(ctx context.Context, args SelectBestCPUArgs) (entity.Part, *errors.ErrService) {
	queueStalePartsUpdate(uc.UpdatePartsUseCase, args.cpus)

	bestCPU := chooseBestPart(args.cpus, partRanking{score: performanceScore})
	if bestCPU == nil {
		return entity.Part{}, nil
	}

	return *bestCPU, nil
}
