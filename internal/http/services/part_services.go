package services

import (
	"github.com/Luzin7/pcideal-be/internal/contracts"
	"github.com/Luzin7/pcideal-be/internal/core/models"
)

type PartService struct {
	PartRepository contracts.PartContract
}

func NewPartService(partRepository contracts.PartContract) *PartService {
	return &PartService{
		PartRepository: partRepository,
	}
}

func (partService *PartService) GetAllParts() ([]*models.Part, error) {
	parts, err := partService.PartRepository.GetAllParts()

	if err != nil {
		return nil, err
	}

	return parts, nil
}
