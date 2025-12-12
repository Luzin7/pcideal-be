package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/util"
)

type GetPartByIDUseCase struct {
	partRepository repository.PartRepository
	updatePartUC   *UpdatePartUseCase
}

func NewGetPartByIDUseCase(partRepository repository.PartRepository, updatePartUC *UpdatePartUseCase) *GetPartByIDUseCase {
	return &GetPartByIDUseCase{
		partRepository: partRepository,
		updatePartUC:   updatePartUC,
	}
}

func (uc *GetPartByIDUseCase) Execute(id string) (*entity.Part, *errors.ErrService) {
	ctx := context.Background()
	part, err := uc.partRepository.GetPartByID(ctx, id)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if part == nil {
		return nil, errors.ErrNotFound("part")
	}

	if util.PartNeedToUpdate(part) {
		uc.updatePartUC.Execute(part.ID.Hex())
	}

	return part, nil
}
