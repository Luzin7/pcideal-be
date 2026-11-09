package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

const (
	DefaultPartsPageSize int64 = 50
	MaxPartsPageSize     int64 = 200
)

type FindAllPartsArgs struct {
	Limit  int64
	Offset int64
}

type FindPartByTypeAndBrandWithMaxPriceArgs struct {
	PartType           entity.PartType
	Brand              string
	MaxPriceCents      int64
	MinPSUWatts        int16
	Socket             string
	MemoryType         string
	IntegratedGraphics bool
}

type PartRepository interface {
	GetAllParts(ctx context.Context, args FindAllPartsArgs) ([]*entity.Part, error)
	GetPartByID(ctx context.Context, id string) (*entity.Part, error)
	UpdatePart(ctx context.Context, partId string, part *entity.Part) error
	FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error)
}
