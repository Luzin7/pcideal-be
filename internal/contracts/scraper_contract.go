package contracts

import (
	"github.com/Luzin7/pcideal-be/infra/dto"
	"github.com/Luzin7/pcideal-be/internal/core/models"
)

type ScraperClient interface {
	ScrapeAllCategories() ([]*models.Part, error)
	ScrapeProduct(productLink string) (*models.Part, error)
	UpdateProducts(urls []*dto.ProductLinkToUpdate) ([]*dto.PartWithID, error)
}
