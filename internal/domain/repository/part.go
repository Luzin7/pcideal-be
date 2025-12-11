package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
)

type FindPartByTypeAndBrandWithMaxPriceArgs struct {
	PartType      entity.PartType
	Brand         string
	MaxPriceCents int64
	MinPSUWatts   int16
	Socket        string
	MemoryType    string
}

type PartRepository interface {
	CreatePart(ctx context.Context, part *entity.Part) error
	GetAllParts(ctx context.Context) ([]*entity.Part, error)
	GetPartByID(ctx context.Context, id string) (*entity.Part, error)
	GetPartByModel(ctx context.Context, model string) (*entity.Part, error)
	UpdatePart(ctx context.Context, partId string, part *entity.Part) error
	UpdateParts(ctx context.Context, urls []dto.ProductLinkToUpdate, store string) error
	FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error)
}
