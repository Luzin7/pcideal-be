package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type GetAllPartsUseCase struct {
	partRepository repository.PartRepository
}

func NewGetAllPartsUseCase(partRepository repository.PartRepository) *GetAllPartsUseCase {
	return &GetAllPartsUseCase{partRepository: partRepository}
}

func (uc *GetAllPartsUseCase) Execute() ([]*entity.Part, *errors.ErrService) {
	ctx := context.Background()
	parts, err := uc.partRepository.GetAllParts(ctx)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if len(parts) == 0 {
		return nil, errors.ErrNotFound("parts")
	}

	return parts, nil
}
