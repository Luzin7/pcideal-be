package presenters

import "github.com/Luzin7/pcideal-be/internal/core/models"

type RecommendedBuildsPresenter struct {
	Builds []RecommendationBuild `json:"builds"`
}

type RecommendationBuild struct {
	BuildType   string     `json:"build_type"`
	Budget      int64      `json:"budget"`
	Description string     `json:"description"`
	Parts       BuildParts `json:"parts"`
	Summary     string     `json:"summary"`
}

type BuildParts struct {
	CPU            *models.Part `json:"cpu"`
	Motherboard    *models.Part `json:"mobo"`
	RAM            *models.Part `json:"ram"`
	PrimaryStorage *models.Part `json:"primary_storage"`
	GPU            *models.Part `json:"gpu"`
	PSU            *models.Part `json:"psu"`
}
