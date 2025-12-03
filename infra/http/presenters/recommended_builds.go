package presenters

import "github.com/Luzin7/pcideal-be/internal/domain/entity"

type RecommendedBuildsPresenter struct {
	Builds []RecommendationBuild `json:"builds"`
}

type RecommendationBuild struct {
	BuildType   string     `json:"build_type"`
	Budget      int64      `json:"budget"`
	BuildValue  int64      `json:"build_value"`
	Description string     `json:"description"`
	Parts       BuildParts `json:"parts"`
	Summary     string     `json:"summary"`
}

type BuildParts struct {
	CPU            *entity.Part `json:"cpu"`
	Motherboard    *entity.Part `json:"mobo"`
	RAM            *entity.Part `json:"ram"`
	PrimaryStorage *entity.Part `json:"primary_storage"`
	GPU            *entity.Part `json:"gpu"`
	PSU            *entity.Part `json:"psu"`
}
