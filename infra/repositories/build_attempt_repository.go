package repositories

import (
	"context"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const maxBuildAttemptsPerWindow = 4

type BuildAttemptRepository struct {
	collection *mongo.Collection
}

func NewBuildAttemptRepository(client *mongo.Database) *BuildAttemptRepository {
	return &BuildAttemptRepository{
		collection: client.Collection("build_attempts"),
	}
}

func (buildattemptRepository *BuildAttemptRepository) ClaimBuildAttempt(ctx context.Context, buildattempt *entity.BuildAttempt) error {
	return claimAttemptSlot(ctx, buildattempt, func(ctx context.Context, attempt *entity.BuildAttempt) error {
		_, err := buildattemptRepository.collection.InsertOne(ctx, attempt)
		return err
	})
}

func (buildattemptRepository *BuildAttemptRepository) ReleaseBuildAttempt(ctx context.Context, buildattempt *entity.BuildAttempt) error {
	_, err := buildattemptRepository.collection.DeleteOne(ctx, buildAttemptReleaseFilter(buildattempt))
	return err
}

func buildAttemptReleaseFilter(attempt *entity.BuildAttempt) bson.M {
	return bson.M{
		"ip":           attempt.IP,
		"window_start": attempt.WindowStart,
		"slot":         attempt.Slot,
	}
}

func claimAttemptSlot(ctx context.Context, attempt *entity.BuildAttempt, insert func(context.Context, *entity.BuildAttempt) error) error {
	windowStart := attempt.CreatedAt.UTC().Truncate(time.Hour)

	for slot := 0; slot < maxBuildAttemptsPerWindow; slot++ {
		attempt.WindowStart = windowStart
		attempt.Slot = slot

		err := insert(ctx, attempt)
		if err == nil {
			return nil
		}
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}

	return repository.ErrBuildAttemptLimitReached
}
