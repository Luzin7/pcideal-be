package services

import (
	"log"
	"time"

	"github.com/Luzin7/pcideal-be/internal/contracts"
	"github.com/Luzin7/pcideal-be/internal/core/models"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type PartService struct {
	PartRepository contracts.PartContract
	ScraperClient  contracts.ScraperClient
}

func NewPartService(partRepository contracts.PartContract, scrapperClient contracts.ScraperClient) *PartService {
	return &PartService{
		PartRepository: partRepository,
		ScraperClient:  scrapperClient,
	}
}

func (partService *PartService) CreatePart(part *models.Part) *errors.ErrService {
	partAlreadyExist, err := partService.PartRepository.GetPartByModel(part.Model)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	if partAlreadyExist != nil {
		return errors.ErrAlreadyExists("part")
	}

	newPart := models.NewPart(
		part.Type,
		part.Brand,
		part.Model,
		part.Specs,
		part.PriceCents,
		part.URL,
		part.Store,
	)

	err = partService.PartRepository.CreatePart(newPart)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	return nil
}

func (partService *PartService) GetAllParts() ([]*models.Part, *errors.ErrService) {
	parts, err := partService.PartRepository.GetAllParts()

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if len(parts) == 0 {
		return nil, errors.ErrNotFound("parts")
	}

	return parts, nil
}

func (partService *PartService) UpdatePart(partId string) *errors.ErrService {
	part, err := partService.PartRepository.GetPartByID(partId)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	if part == nil {
		return errors.ErrNotFound("part")
	}

	scrapedPart, err := partService.ScraperClient.ScrapeProduct(part.URL)

	if err != nil {
		return errors.ErrInternalServerError()
	}

	part.Specs = scrapedPart.Specs
	part.PriceCents = scrapedPart.PriceCents
	part.UpdatedAt = time.Now()

	err = partService.PartRepository.UpdatePart(partId, part)

	if err != nil {
		log.Printf("Error updating part: %v", err)
		return errors.ErrInternalServerError()
	}

	return nil
}

func (partService *PartService) GetPartByID(id string) (*models.Part, *errors.ErrService) {
	part, err := partService.PartRepository.GetPartByID(id)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if part == nil {
		return nil, errors.ErrNotFound("part")
	}

	if time.Since(part.UpdatedAt) >= 2*time.Hour {
		go func(partID string) {
			if err := partService.UpdatePart(partID); err != nil {
				log.Printf("async update error for part %s: %v", partID, err)
			}
		}(part.ID.Hex())
	}

	return part, nil
}
func (partService *PartService) GetPartByModel(model string) (*models.Part, *errors.ErrService) {
	part, err := partService.PartRepository.GetPartByModel(model)

	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	if part == nil {
		return nil, errors.ErrNotFound("part")
	}

	if time.Since(part.UpdatedAt) >= 2*time.Hour {
		go func(partID string) {
			if err := partService.UpdatePart(partID); err != nil {
				log.Printf("async update error for part %s: %v", partID, err)
			}
		}(part.ID.Hex())
	}

	return part, nil
}
