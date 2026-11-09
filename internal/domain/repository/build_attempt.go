package repository

import (
	"context"
	stderrors "errors"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

var ErrBuildAttemptLimitReached = stderrors.New("build attempt limit reached")

type BuildAttemptRepository interface {
	ClaimBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) error
	ReleaseBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) error
}
