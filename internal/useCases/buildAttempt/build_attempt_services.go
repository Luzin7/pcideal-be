package buildAttempt

import (
	"context"
	"log"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"github.com/Luzin7/pcideal-be/internal/errors"
)

type BuildAttemptService struct {
	BuildAttemptRepository repository.BuildAttemptRepository
}

func NewBuildAttemptService(buildattemptRepository repository.BuildAttemptRepository) *BuildAttemptService {
	return &BuildAttemptService{
		BuildAttemptRepository: buildattemptRepository,
	}
}

func (s *BuildAttemptService) CreateBuildAttempt(buildAttempt *entity.BuildAttempt) *errors.ErrService {
	if buildAttempt.IP == "" {
		return errors.ErrMissingIP()
	}

	newBuildAttempt := entity.NewBuildAttempt(
		buildAttempt.IP, buildAttempt.Goal, buildAttempt.Budget, buildAttempt.CPUPref, buildAttempt.GPUPref,
	)

	ctx := context.Background()
	err := s.BuildAttemptRepository.CreateBuildAttempt(ctx, newBuildAttempt)
	if err == nil {
		return nil
	}
	return errors.ErrInternalServerError()
}
func (s *BuildAttemptService) CountBuildAttemptsByIP(ip string, since time.Time) (int, *errors.ErrService) {
	if ip == "" {
		return 0, errors.ErrMissingIP()
	}
	if since.After(time.Now()) {
		return 0, errors.ErrInvalidSince()
	}
	ctx := context.Background()
	count, err := s.BuildAttemptRepository.CountBuildAttemptsByIP(ctx, ip, since)
	log.Print(count)
	if err != nil {
		return 0, errors.ErrInternalServerError()
	}
	return count, nil
}

func (s *BuildAttemptService) GetBuildAttemptsByIP(ip string, since time.Time) ([]*entity.BuildAttempt, *errors.ErrService) {
	if ip == "" {
		return nil, errors.ErrMissingIP()
	}
	if since.After(time.Now()) {
		return nil, errors.ErrInvalidSince()
	}
	ctx := context.Background()
	attempts, err := s.BuildAttemptRepository.GetBuildAttemptsByIP(ctx, ip, since)
	if err != nil {
		return nil, errors.ErrInternalServerError()
	}
	return attempts, nil
}
