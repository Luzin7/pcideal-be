package part

import (
	"github.com/Luzin7/pcideal-be/infra/http/presenters"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type GenerateBuildRecommendationsUseCase struct {
	partRepository repository.PartRepository
	scraperClient  repository.ScraperClientRepository
	googleAIClient repository.GoogleAIRepository
	updatePartsUC  *UpdatePartsUseCase
	// TODO: Adicionar PartMatchingService quando implementado
}

func NewGenerateBuildRecommendationsUseCase(
	partRepository repository.PartRepository,
	scraperClient repository.ScraperClientRepository,
	googleAIClient repository.GoogleAIRepository,
	updatePartsUC *UpdatePartsUseCase,
) *GenerateBuildRecommendationsUseCase {
	return &GenerateBuildRecommendationsUseCase{
		partRepository: partRepository,
		scraperClient:  scraperClient,
		googleAIClient: googleAIClient,
		updatePartsUC:  updatePartsUC,
	}
}

func (uc *GenerateBuildRecommendationsUseCase) Execute(usageType string, cpuPreference string, gpuPreference string, budget int64) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	// TODO: Implementar PartMatchingService para FindParts e FindBestMatch
	// TODO: Implementar função ValidateCPUAndMotherboard para validar compatibilidade entre CPU e Motherboard

	return nil, errors.New("Build recommendations not yet implemented", 501)
}
