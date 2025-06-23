package services

import (
	"log"
	"strings"
	"time"

	"github.com/Luzin7/pcideal-be/internal/contracts"
	"github.com/Luzin7/pcideal-be/internal/core/models"
	"github.com/Luzin7/pcideal-be/internal/domain/validation"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/Luzin7/pcideal-be/internal/http/presenters"
)

type PartService struct {
	PartRepository      contracts.PartContract
	ScraperClient       contracts.ScraperClient
	GoogleAIClient      contracts.GoogleAIContract
	PartMatchingService contracts.PartMatchingContract
}

func NewPartService(partRepository contracts.PartContract, scrapperClient contracts.ScraperClient, googleAIContract contracts.GoogleAIContract, partMatchingService contracts.PartMatchingContract) *PartService {
	return &PartService{
		PartRepository:      partRepository,
		ScraperClient:       scrapperClient,
		GoogleAIClient:      googleAIContract,
		PartMatchingService: partMatchingService,
	}
}

func (ps *PartService) updatePartIfNeeded(part *models.Part) {
	if part != nil && time.Since(part.UpdatedAt) >= 2*time.Hour {
		go func(partID string) {
			if err := ps.UpdatePart(partID); err != nil {
				log.Printf("async update error for part %s: %v", partID, err)
			}
		}(part.ID.Hex())
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
		part.AffiliatedURL,
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
		log.Printf("Error scraping product for part %s: %v", partId, err)
		return errors.ErrScrapingFailed(partId)
	}

	part.Brand = scrapedPart.Brand
	part.Model = scrapedPart.Model
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

	partService.updatePartIfNeeded(part)

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

	partService.updatePartIfNeeded(part)

	return part, nil
}

func (partService *PartService) GenerateBuildRecomendations(usageType string, cpuPreference string, gpuPreference string, budget int64) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	prompt, err := partService.GoogleAIClient.BuildComputerPrompt(usageType, cpuPreference, gpuPreference, budget)
	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	aiBuildResponse, err := partService.GoogleAIClient.GenerateBuilds(prompt)
	if err != nil {
		return nil, errors.ErrInternalServerError()
	}

	recommendationBuilds := make([]presenters.RecommendationBuild, 0)

	for i := range aiBuildResponse.Builds {
		cpu := aiBuildResponse.Builds[i].Parts.CPU.Model
		cpuBrand := aiBuildResponse.Builds[i].Parts.CPU.Brand
		mobo := aiBuildResponse.Builds[i].Parts.Motherboard.Model
		moboBrand := aiBuildResponse.Builds[i].Parts.Motherboard.Brand
		ram := aiBuildResponse.Builds[i].Parts.RAM.Model
		ramBrand := aiBuildResponse.Builds[i].Parts.RAM.Brand
		gpu := aiBuildResponse.Builds[i].Parts.GPU.Model
		gpuBrand := aiBuildResponse.Builds[i].Parts.GPU.Brand
		primaryStorage := aiBuildResponse.Builds[i].Parts.PrimaryStorage.Model
		primaryStorageBrand := aiBuildResponse.Builds[i].Parts.PrimaryStorage.Brand
		psu := aiBuildResponse.Builds[i].Parts.PSU.Model
		psuBrand := aiBuildResponse.Builds[i].Parts.PSU.Brand

		cpuParts, err := partService.PartMatchingService.FindParts(cpu, "cpu", cpuBrand)
		if err != nil || len(cpuParts) == 0 {
			log.Print(err)
			continue
		}

		moboParts, err := partService.PartMatchingService.FindParts(mobo, "motherboard", moboBrand)
		if err != nil || len(moboParts) == 0 {
			log.Print(err)
			continue
		}

		cpuFoundByBestMatch := partService.PartMatchingService.FindBestMatch(cpu, cpuParts)
		if cpuFoundByBestMatch == nil {
			log.Print(err)
			continue
		}
		partService.updatePartIfNeeded(cpuFoundByBestMatch)

		moboFoundByBestMatch := partService.PartMatchingService.FindBestMatch(mobo, moboParts)
		if moboFoundByBestMatch == nil {
			log.Print(err)
			continue
		}
		partService.updatePartIfNeeded(moboFoundByBestMatch)

		isBuildValid := validation.ValidateCPUAndMotherboard(cpuFoundByBestMatch, moboFoundByBestMatch)

		if !isBuildValid {
			continue
		}

		ramParts, err := partService.PartMatchingService.FindParts(ram, "ram", ramBrand)
		if err != nil || len(ramParts) == 0 {
			log.Print(err)
			continue
		}

		primaryStorageParts, err := partService.PartMatchingService.FindParts(primaryStorage, "ssd", primaryStorageBrand)
		if err != nil || len(primaryStorageParts) == 0 {
			log.Print(err)
			continue
		}

		psuParts, err := partService.PartMatchingService.FindParts(psu, "psu", psuBrand)
		if err != nil || len(psuParts) == 0 {
			log.Print(err)
			continue
		}

		var gpuFoundByBestMatch *models.Part

		if strings.ToLower(gpu) != "integrado" {
			gpuParts, err := partService.PartMatchingService.FindParts(gpu, "gpu", gpuBrand)
			if err != nil || len(gpuParts) == 0 {
				log.Print(err)
				continue
			}
			gpuFoundByBestMatch = partService.PartMatchingService.FindBestMatch(gpu, gpuParts)
			if gpuFoundByBestMatch == nil {
				log.Print(err)
				continue
			}
			partService.updatePartIfNeeded(gpuFoundByBestMatch)
		} else {
			gpuFoundByBestMatch = &models.Part{
				Brand: gpuBrand,
				Model: "Integrada",
				URL:   cpuFoundByBestMatch.URL,
			}
		}

		ramFoundByBestMatch := partService.PartMatchingService.FindBestMatch(ram, ramParts)
		if ramFoundByBestMatch == nil {
			log.Print(err)
			continue
		}
		partService.updatePartIfNeeded(ramFoundByBestMatch)

		primaryStorageFoundByBestMatch := partService.PartMatchingService.FindBestMatch(primaryStorage, primaryStorageParts)
		if primaryStorageFoundByBestMatch == nil {
			log.Print(err)
			continue
		}
		partService.updatePartIfNeeded(primaryStorageFoundByBestMatch)

		psuFoundByBestMatch := partService.PartMatchingService.FindBestMatch(psu, psuParts)
		if psuFoundByBestMatch == nil {
			log.Print(err)
			continue
		}
		partService.updatePartIfNeeded(psuFoundByBestMatch)

		buildValue := cpuFoundByBestMatch.PriceCents +
			moboFoundByBestMatch.PriceCents +
			ramFoundByBestMatch.PriceCents +
			primaryStorageFoundByBestMatch.PriceCents +
			psuFoundByBestMatch.PriceCents

		if gpuFoundByBestMatch.Model != "Integrado" {
			buildValue += gpuFoundByBestMatch.PriceCents
		}

		recommendationBuild := presenters.RecommendationBuild{
			BuildType:   aiBuildResponse.Builds[i].BuildType,
			Budget:      aiBuildResponse.Builds[i].Budget,
			BuildValue:  buildValue,
			Description: aiBuildResponse.Builds[i].Description,
			Summary:     aiBuildResponse.Builds[i].Summary,
			Parts: presenters.BuildParts{
				CPU:            cpuFoundByBestMatch,
				Motherboard:    moboFoundByBestMatch,
				RAM:            ramFoundByBestMatch,
				GPU:            gpuFoundByBestMatch,
				PrimaryStorage: primaryStorageFoundByBestMatch,
				PSU:            psuFoundByBestMatch,
			},
		}

		recommendationBuilds = append(recommendationBuilds, recommendationBuild)

	}

	if len(recommendationBuilds) == 0 {
		return nil, errors.ErrNotFound("compatible build")
	}

	return &presenters.RecommendedBuildsPresenter{
		Builds: recommendationBuilds,
	}, nil
}
