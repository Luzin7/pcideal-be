package contracts

import "github.com/Luzin7/pcideal-be/internal/core/models"

type PartContract interface {
	GetAllParts() ([]*models.Part, error)
	GetPartByID(id string) (*models.Part, error)
	GetPartByName(name string) (*models.Part, error)
}
