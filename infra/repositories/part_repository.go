package repositories

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (partRepository *PartRepository) CreatePart(part *models.Part) error {
	ctx := context.TODO()

	_, err := partRepository.collection.InsertOne(ctx, part)

	if err != nil {
		return err
	}

	return nil
}

func (partRepository *PartRepository) UpdatePart(partId string, part *models.Part) error {
	ctx := context.TODO()
	log.Printf("Updating part with ID: %s", partId)
	log.Printf("Part details: %+v", part)

	objID, err := primitive.ObjectIDFromHex(partId)

	if err != nil {
		log.Printf("Error converting part ID to ObjectID: %v", err)
		return err
	}

	_, err = partRepository.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
		"$set": bson.M{
			"type":        part.Type,
			"brand":       part.Brand,
			"model":       part.Model,
			"specs":       part.Specs,
			"price_cents": part.PriceCents,
			"url":         part.URL,
			"store":       part.Store,
			"updated_at":  part.UpdatedAt,
		},
	})

	if err != nil {
		log.Printf("Error updating part: %v", err)
		return err
	}

	return nil
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

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = partRepository.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&part)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &part, nil
}

func (partRepository *PartRepository) GetPartByModel(model string) (*models.Part, error) {
	ctx := context.TODO()

	var part models.Part

	err := partRepository.collection.FindOne(ctx, bson.M{"model": model}).Decode(&part)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &part, nil
}
