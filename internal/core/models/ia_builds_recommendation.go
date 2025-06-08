package models

type AIBuildResponse struct {
	Builds []AIBuild `json:"builds"`
}

type AIBuild struct {
	BuildType   string  `json:"build_type"`
	Budget      int64   `json:"budget"`
	Description string  `json:"description"`
	Parts       AIParts `json:"parts"`
	Summary     string  `json:"summary"`
}

type AIParts struct {
	CPU            AIPartInfo `json:"cpu"`
	Motherboard    AIPartInfo `json:"mobo"`
	RAM            AIPartInfo `json:"ram"`
	PrimaryStorage AIPartInfo `json:"primary_storage"`
	GPU            AIPartInfo `json:"gpu"`
	PSU            AIPartInfo `json:"psu"`
}

type AIPartInfo struct {
	Model string `json:"model"`
	Brand string `json:"brand"`
}
