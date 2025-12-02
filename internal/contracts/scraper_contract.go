package contracts

import (
	"github.com/Luzin7/pcideal-be/infra/dto"
	"github.com/Luzin7/pcideal-be/internal/core/models"
)

type ScraperClient interface {
	ScrapeAllCategories(storeName string) ([]*models.Part, error)
	ScrapeProduct(productLink string, storeName string) (*models.Part, error)
	UpdateProducts(urls []*dto.ProductLinkToUpdate, storeName string) ([]*dto.PartWithID, error)
}
