package part

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type CreatePartUseCase struct {
	partRepository repository.PartRepository
}

func NewCreatePartUseCase(partRepository repository.PartRepository) *CreatePartUseCase {
	return &CreatePartUseCase{partRepository: partRepository}
}

func (uc *CreatePartUseCase) Execute(part *entity.Part) *errors.ErrService {
	ctx := context.Background()
	partAlreadyExist, err := uc.partRepository.GetPartByModel(ctx, part.Model)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	if partAlreadyExist != nil {
		return errors.ErrAlreadyExists("part")
	}

	newPart := entity.NewPart(
		part.Type,
		part.Brand,
		part.Model,
		part.URL,
		part.Store,
		part.AffiliatedURL,
		part.PriceCents,
		part.Specs,
	)

	err = uc.partRepository.CreatePart(ctx, newPart)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	return nil
}
