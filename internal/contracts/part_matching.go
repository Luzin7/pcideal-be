package contracts

import "github.com/Luzin7/pcideal-be/internal/core/models"

type PartMatchingContract interface {
	FindParts(productName string, productType string, productBrand string) ([]models.Part, error)
	FindBestMatch(targetName string, parts []models.Part) *models.Part
}
