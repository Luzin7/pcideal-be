package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
)

type ScraperClientRepository interface {
	ScrapeAllCategories(ctx context.Context, storeName string) ([]*entity.Part, error)
	ScrapeProduct(ctx context.Context, productLink string, storeName string) (*entity.Part, error)
	UpdateProducts(ctx context.Context, urls []*dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error)
}
