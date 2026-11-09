package buildAttempt

import (
	"context"
	stderrors "errors"

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

func (s *BuildAttemptService) ClaimBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) *errors.ErrService {
	if buildAttempt.IP == "" {
		return errors.ErrMissingIP()
	}

	err := s.BuildAttemptRepository.ClaimBuildAttempt(ctx, buildAttempt)
	if err == nil {
		return nil
	}
	if stderrors.Is(err, repository.ErrBuildAttemptLimitReached) {
		return errors.ErrBuildAttemptLimit()
	}
	return errors.ErrInternalServerError()
}

func (s *BuildAttemptService) ReleaseBuildAttempt(ctx context.Context, buildAttempt *entity.BuildAttempt) *errors.ErrService {
	if err := s.BuildAttemptRepository.ReleaseBuildAttempt(ctx, buildAttempt); err != nil {
		return errors.ErrInternalServerError()
	}

	return nil
}
