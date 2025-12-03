package repository

import (
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
)

type BuildAttemptRepository interface {
	CreateBuildAttempt(buildAttempt *entity.BuildAttempt) error
	CountBuildAttemptsByIP(ip string, since time.Time) (int, error)
	GetBuildAttemptsByIP(ip string, since time.Time) ([]*entity.BuildAttempt, error)
}
