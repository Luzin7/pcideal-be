package buildAttempt

import (
	"context"
	"errors"
	"testing"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
)

type stubBuildAttemptRepository struct {
	claimErr   error
	releaseErr error
}

func (s *stubBuildAttemptRepository) ClaimBuildAttempt(context.Context, *entity.BuildAttempt) error {
	return s.claimErr
}

func (s *stubBuildAttemptRepository) ReleaseBuildAttempt(context.Context, *entity.BuildAttempt) error {
	return s.releaseErr
}

func TestClaimBuildAttemptMapsLimitTo429(t *testing.T) {
	service := NewBuildAttemptService(&stubBuildAttemptRepository{claimErr: repository.ErrBuildAttemptLimitReached})

	err := service.ClaimBuildAttempt(context.Background(), &entity.BuildAttempt{IP: "203.0.113.7"})

	if err == nil || err.StatusCode != 429 {
		t.Fatalf("expected 429, got %#v", err)
	}
}

func TestClaimBuildAttemptMapsUnknownErrorTo500(t *testing.T) {
	service := NewBuildAttemptService(&stubBuildAttemptRepository{claimErr: errors.New("boom")})

	err := service.ClaimBuildAttempt(context.Background(), &entity.BuildAttempt{IP: "203.0.113.7"})

	if err == nil || err.StatusCode != 500 {
		t.Fatalf("expected 500, got %#v", err)
	}
}

func TestClaimBuildAttemptRejectsMissingIP(t *testing.T) {
	service := NewBuildAttemptService(&stubBuildAttemptRepository{})

	err := service.ClaimBuildAttempt(context.Background(), &entity.BuildAttempt{})

	if err == nil || err.StatusCode != 400 {
		t.Fatalf("expected 400, got %#v", err)
	}
}

func TestReleaseBuildAttemptMapsUnknownErrorTo500(t *testing.T) {
	service := NewBuildAttemptService(&stubBuildAttemptRepository{releaseErr: errors.New("boom")})

	err := service.ReleaseBuildAttempt(context.Background(), &entity.BuildAttempt{IP: "203.0.113.7"})

	if err == nil || err.StatusCode != 500 {
		t.Fatalf("expected 500, got %#v", err)
	}
}
