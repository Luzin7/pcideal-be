package contracts

type ValidateBuildContract interface {
	ValidateCPUAndMotherboard(cpu string, mobo string) bool
}
