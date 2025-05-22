package repositories

import (
	"context"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PartRepository struct {
	collection *mongo.Collection
}

func NewPartRepository(client *mongo.Database) *PartRepository {
	collection := client.Collection("parts")

	return &PartRepository{
		collection: collection,
	}
}

func (partRepository *PartRepository) GetAllParts() ([]*models.Part, error) {
	ctx := context.TODO()

	cursor, err := partRepository.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var parts []*models.Part

	for cursor.Next(ctx) {
		var part models.Part

		if err := cursor.Decode(&part); err != nil {
			return nil, err
		}

		parts = append(parts, &part)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}
