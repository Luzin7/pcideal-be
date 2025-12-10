package part

import (
	"context"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type GetPartByModelUseCase struct {
	partRepository repository.PartRepository
	updatePartUC   *UpdatePartUseCase
}

func NewGetPartByModelUseCase(partRepository repository.PartRepository, updatePartUC *UpdatePartUseCase) *GetPartByModelUseCase {
	return &GetPartByModelUseCase{
		partRepository: partRepository,
		updatePartUC:   updatePartUC,
	}
}

func (uc *GetPartByModelUseCase) Execute(model string) (*entity.Part, *errors.ErrService) {
	ctx := context.Background()
	part, err := uc.partRepository.GetPartByModel(ctx, model)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if part == nil {
		return nil, errors.ErrNotFound("part")
	}

	if uc.needToUpdate(part) {
		uc.updatePartUC.Execute(part.ID.Hex())
	}

	return part, nil
}

func (uc *GetPartByModelUseCase) needToUpdate(part *entity.Part) bool {
	if part != nil && time.Since(part.UpdatedAt) >= 2*time.Hour {
		return true
	}
	return false
}
