package repository

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type GoogleAIRepository interface {
	GeneratePcBuildAnalysis(ctx context.Context, part *entity.Part) (string, error)
}
