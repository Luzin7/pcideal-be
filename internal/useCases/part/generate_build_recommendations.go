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
	// TODO: Implementar integração com GoogleAIClient quando os métodos BuildComputerPrompt e GenerateBuilds forem adicionados à interface
	// TODO: Implementar PartMatchingService para FindParts e FindBestMatch
	// TODO: Implementar função ValidateCPUAndMotherboard para validar compatibilidade entre CPU e Motherboard

	// Lógica completa será implementada após criar os serviços necessários:
	// 1. GoogleAIClient.BuildComputerPrompt() - Gerar prompt baseado nos parâmetros
	// 2. GoogleAIClient.GenerateBuilds() - Gerar recomendações de build via IA
	// 3. PartMatchingService.FindParts() - Buscar peças compatíveis no banco
	// 4. PartMatchingService.FindBestMatch() - Encontrar melhor match de peças
	// 5. validation.ValidateCPUAndMotherboard() - Validar compatibilidade de socket
	// 6. Atualizar preços em background via UpdatePartsUseCase

	return nil, errors.New("Build recommendations not yet implemented", 501)
}
