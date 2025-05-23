package contracts

import "github.com/Luzin7/pcideal-be/internal/core/models"

type ScraperClient interface {
	ScrapeAllCategories() ([]*models.Part, error)
	ScrapeProduct(productLink string) (*models.Part, error)
}
