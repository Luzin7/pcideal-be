package dto

import (
	"regexp"

	"github.com/Luzin7/pcideal-be/internal/errors"
)

const maxPreferenceLength = 50

var preferencePattern = regexp.MustCompile(`^[a-zA-Z0-9 ._+-]*$`)

type GenerateBuildRecommendationsDTO struct {
	UsageType     string `json:"usage_type"`
	CpuPreference string `json:"cpu_preference"`
	GpuPreference string `json:"gpu_preference"`
	BudgetCents   int64  `json:"budget_cents" binding:"required,gt=0"`
}

func (dto GenerateBuildRecommendationsDTO) Validate() *errors.ErrService {
	if !isAllowedUsageType(dto.UsageType) {
		return errors.ErrInvalidUsageType(dto.UsageType)
	}

	if dto.BudgetCents <= 0 {
		return errors.ErrInvalidBudget()
	}

	if err := validatePreference("cpu_preference", dto.CpuPreference); err != nil {
		return err
	}

	return validatePreference("gpu_preference", dto.GpuPreference)
}

func isAllowedUsageType(usageType string) bool {
	switch usageType {
	case "GAMING", "WORK", "CONTENT_CREATOR":
		return true
	default:
		return false
	}
}

func validatePreference(field string, preference string) *errors.ErrService {
	if preference == "" {
		return nil
	}

	if len(preference) > maxPreferenceLength || !preferencePattern.MatchString(preference) {
		return errors.ErrInvalidPreference(field)
	}

	return nil
}
