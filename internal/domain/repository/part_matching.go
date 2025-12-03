package repository

import "github.com/Luzin7/pcideal-be/internal/domain/entity"

type PartMatchingRepository interface {
	FindParts(productName string, productType string, productBrand string) ([]entity.Part, error)
	FindBestMatch(targetName string, parts []entity.Part) *entity.Part
}
