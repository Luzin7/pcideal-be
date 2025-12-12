package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
)

type GoogleAIRepository interface {
	GeneratePcBuildAnalysis(ctx context.Context, build *presenters.RecommendationBuild) (string, error)
}
