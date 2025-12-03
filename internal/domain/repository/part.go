package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type PartRepository interface {
	CreatePart(ctx context.Context, part *entity.Part) error
	GetAllParts(ctx context.Context) ([]*entity.Part, error)
	GetPartByID(ctx context.Context, id string) (*entity.Part, error)
	GetPartByModel(ctx context.Context, model string) (*entity.Part, error)
	UpdatePart(ctx context.Context, partId string, part *entity.Part) error
}
