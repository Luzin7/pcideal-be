package repository

import "github.com/Luzin7/pcideal-be/internal/domain/entity"

type GoogleAIRepository interface {
	BuildComputerPrompt(usageType string, cpuPreference string, gpuPreference string, budget int64) (string, error)
	GenerateBuilds(prompt string) (*entity.AIBuildResponse, error)
}
