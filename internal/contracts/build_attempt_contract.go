package contracts

import (
	"time"

	"github.com/Luzin7/pcideal-be/internal/core/models"
)

type BuildAttemptContract interface {
	CreateBuildAttempt(buildAttempt *models.BuildAttempt) error
	CountBuildAttemptsByIP(ip string, since time.Time) (int, error)
	GetBuildAttemptsByIP(ip string, since time.Time) ([]*models.BuildAttempt, error)
}
