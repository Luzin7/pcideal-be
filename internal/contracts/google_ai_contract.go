package contracts

type GoogleAIContract interface {
	BuildComputerPrompt(usageType string, cpuPreference string, gpuPreference string, budget int64) (string, error)
	GenerateBuilds(prompt string) (map[string]interface{}, error)
}
