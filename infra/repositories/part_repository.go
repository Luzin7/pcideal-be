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

func (partRepository *PartRepository) GetPartByID(id string) (*models.Part, error) {
	ctx := context.TODO()

	var part models.Part

	err := partRepository.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&part)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &part, nil
}

func (partRepository *PartRepository) GetPartByName(name string) (*models.Part, error) {
	ctx := context.TODO()

	var part models.Part

	err := partRepository.collection.FindOne(ctx, bson.M{"name": name}).Decode(&part)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &part, nil
}
