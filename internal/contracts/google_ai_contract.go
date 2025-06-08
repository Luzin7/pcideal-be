package contracts

import "github.com/Luzin7/pcideal-be/internal/core/models"

type GoogleAIContract interface {
	BuildComputerPrompt(usageType string, cpuPreference string, gpuPreference string, budget int64) (string, error)
	GenerateBuilds(prompt string) (*models.AIBuildResponse, error)
}
