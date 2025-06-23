package repositories

import (
	"context"
	"time"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BuildAttemptRepository struct {
	collection *mongo.Collection
}

func NewBuildAttemptRepository(client *mongo.Database) *BuildAttemptRepository {
	collection := client.Collection("build_attempts")

	return &BuildAttemptRepository{
		collection: collection,
	}
}

func (buildattemptRepository *BuildAttemptRepository) CreateBuildAttempt(buildattempt *models.BuildAttempt) error {
	ctx := context.TODO()

	_, err := buildattemptRepository.collection.InsertOne(ctx, buildattempt)

	if err != nil {
		return err
	}

	return nil
}

func (buildattemptRepository *BuildAttemptRepository) CountBuildAttemptsByIP(ip string, since time.Time) (int, error) {
	ctx := context.TODO()

	filter := bson.M{
		"ip":         ip,
		"created_at": bson.M{"$gte": since},
	}

	count, err := buildattemptRepository.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (buildattemptRepository *BuildAttemptRepository) GetBuildAttemptsByIP(ip string, since time.Time) ([]*models.BuildAttempt, error) {
	ctx := context.TODO()

	filter := bson.M{
		"ip":         ip,
		"created_at": bson.M{"$gte": since},
	}

	cursor, err := buildattemptRepository.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var buildattempts []*models.BuildAttempt
	for cursor.Next(ctx) {
		var buildattempt models.BuildAttempt
		if err := cursor.Decode(&buildattempt); err != nil {
			return nil, err
		}
		buildattempts = append(buildattempts, &buildattempt)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return buildattempts, nil
}
