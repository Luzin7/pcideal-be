package contracts

import "github.com/Luzin7/pcideal-be/internal/core/models"

type PartContract interface {
	CreatePart(part *models.Part) error
	GetAllParts() ([]*models.Part, error)
	GetPartByID(id string) (*models.Part, error)
	GetPartByModel(model string) (*models.Part, error)
	UpdatePart(partId string, part *models.Part) error
}
