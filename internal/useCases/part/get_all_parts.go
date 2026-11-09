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

func (uc *GetAllPartsUseCase) Execute(ctx context.Context, args repository.FindAllPartsArgs) ([]*entity.Part, *errors.ErrService) {
	parts, err := uc.partRepository.GetAllParts(ctx, args)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if parts == nil {
		return []*entity.Part{}, nil
	}

	return parts, nil
}
