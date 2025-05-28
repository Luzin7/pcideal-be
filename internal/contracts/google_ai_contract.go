package contracts

type GoogleAIContract interface {
	GenerateBuilds(prompt string) (string, error)
}
