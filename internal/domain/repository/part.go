package repository

import "github.com/Luzin7/pcideal-be/internal/domain/entity"

type PartRepository interface {
	CreatePart(part *entity.Part) error
	GetAllParts() ([]*entity.Part, error)
	GetPartByID(id string) (*entity.Part, error)
	GetPartByModel(model string) (*entity.Part, error)
	UpdatePart(partId string, part *entity.Part) error
}
