package dto

type GenerateBuildRecommendationsDTO struct {
	UsageType     string `json:"usage_type"`
	CpuPreference string `json:"cpu_preference"`
	GpuPreference string `json:"gpu_preference"`
	Budget        int64  `json:"budget"`
}
