package entity

type buildType string

const (
	BuildTypeEconomic    buildType = "ECONOMIC"
	BuildTypeBalanced    buildType = "BALANCED"
	BuildTypePerformance buildType = "PERFORMANCE"
)

type Build struct {
	BuildType   buildType
	Budget      int64
	BuildValue  int64
	Description string
	Parts       BuildParts
}

type BuildParts struct {
	CPU            *Part
	Motherboard    *Part
	RAM            *Part
	PrimaryStorage *Part
	GPU            *Part
	PSU            *Part
}

func NewBuild(buildType buildType, budget int64, buildValue int64, description string, parts BuildParts) *Build {
	return &Build{
		BuildType:   buildType,
		Budget:      budget,
		BuildValue:  buildValue,
		Description: description,
		Parts:       parts,
	}
}
