package services

import (
	"errors"
	"log"
	"time"

	"github.com/Luzin7/pcideal-be/internal/contracts"
	"github.com/Luzin7/pcideal-be/internal/core/models"
	"github.com/Luzin7/pcideal-be/pkg/times"
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

func (partService *PartService) CreatePart(part *models.Part) error {
	partAlreadyExist, err := partService.PartRepository.GetPartByModel(part.Model)

	if err != nil {
		return err
	}

	if partAlreadyExist != nil {
		return errors.New("part already exists")
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
		return err
	}

	return nil
}

func (partService *PartService) GetAllParts() ([]*models.Part, error) {
	parts, err := partService.PartRepository.GetAllParts()

	if err != nil {
		return nil, err
	}

	return parts, nil
}

func (partService *PartService) UpdatePart(partId string) error {
	part, err := partService.PartRepository.GetPartByID(partId)

	if err != nil {
		return errors.New("part not found")
	}

	scrapedPart, err := partService.ScraperClient.ScrapeProduct(part.URL)

	if err != nil {
		return errors.New("error scraping product")
	}

	part.Specs = scrapedPart.Specs
	part.PriceCents = scrapedPart.PriceCents
	part.UpdatedAt = time.Now()

	err = partService.PartRepository.UpdatePart(partId, part)

	if err != nil {
		return errors.New("error updating part")
	}

	return nil
}

func (partService *PartService) GetPartByID(id string) (*models.Part, error) {
	part, err := partService.PartRepository.GetPartByID(id)

	if err != nil {
		return nil, errors.New("part not found")
	}

	if part.UpdatedAt.Unix()-time.Now().Unix() >= times.OneHourInSeconds*2 {
		go func() {
			err := partService.UpdatePart(part.ID)
			if err != nil {
				log.Printf("error updating part: %v", err)
			}
		}()
	}

	return part, nil
}
func (partService *PartService) GetPartByModel(model string) (*models.Part, error) {
	part, err := partService.PartRepository.GetPartByModel(model)

	if err != nil {
		return nil, err
	}

	if part.UpdatedAt.Unix()-time.Now().Unix() >= times.OneHourInSeconds*2 {
		go func() {
			err := partService.UpdatePart(part.ID)
			if err != nil {
				log.Printf("error updating part: %v", err)
			}
		}()
	}

	return part, nil
}
