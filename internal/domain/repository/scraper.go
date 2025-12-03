package repository

import (
	"github.com/Luzin7/pcideal-be/infra/dto"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type ScraperClientRepository interface {
	ScrapeAllCategories(storeName string) ([]*entity.Part, error)
	ScrapeProduct(productLink string, storeName string) (*entity.Part, error)
	UpdateProducts(urls []*dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error)
}
