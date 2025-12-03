package part

import (
	"context"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
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

	if uc.needToUpdate(part) {
		uc.updatePartUC.Execute(part.ID.Hex())
	}

	return part, nil
}

func (uc *GetPartByIDUseCase) needToUpdate(part *entity.Part) bool {
	if part != nil && time.Since(part.UpdatedAt) >= 2*time.Hour {
		return true
	}
	return false
}
