package repository

import (
	"context"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type BuildAttemptRepository interface {
	CreateBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) error
	CountBuildAttemptsByIP(ctx context.Context, ip string, since time.Time) (int, error)
	GetBuildAttemptsByIP(ctx context.Context, ip string, since time.Time) ([]*entity.BuildAttempt, error)
}
