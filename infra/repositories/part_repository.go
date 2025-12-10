package repositories

import (
	"context"
	"log"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (partRepository *PartRepository) CreatePart(ctx context.Context, part *entity.Part) error {

	_, err := partRepository.collection.InsertOne(ctx, part)

	if err != nil {
		return err
	}

	return nil
}

func (partRepository *PartRepository) UpdatePart(ctx context.Context, partId string, part *entity.Part) error {
	log.Printf("Updating part with ID: %s", partId)
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

func (partRepository *PartRepository) GetAllParts(ctx context.Context) ([]*entity.Part, error) {

	cursor, err := partRepository.collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var parts []*entity.Part

	for cursor.Next(ctx) {
		part := new(entity.Part)

		if err := cursor.Decode(part); err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}

func (partRepository *PartRepository) GetPartByID(ctx context.Context, id string) (*entity.Part, error) {

	var part entity.Part

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

func (partRepository *PartRepository) GetPartByModel(ctx context.Context, model string) (*entity.Part, error) {

	var part entity.Part

	err := partRepository.collection.FindOne(ctx, bson.M{"model": model}).Decode(&part)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &part, nil
}

func (partRepository *PartRepository) FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args repository.FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error) {
	filter := bson.M{
		"type":        args.PartType,
		"price_cents": bson.M{"$lte": args.MaxPriceCents},
	}

	if args.Brand != "" {
		filter["brand"] = bson.M{"$regex": args.Brand, "$options": "i"}
	}
	if args.PartType == "CPU" {
		filter["specs.socket"] = bson.M{"$exists": true, "$ne": ""}
	}
	if args.PartType == "PSU" && args.MinPSUWatts > 0 {
		filter["specs.wattage"] = bson.M{"$gte": args.MinPSUWatts}
	}
	if args.PartType == "MOBO" && args.Socket != "" {
		filter["specs.socket"] = args.Socket
	}

	// Ordenar baseado no tipo de peça para otimizar a seleção
	var sortField string
	switch args.PartType {
	case "PSU":
		sortField = "specs.efficiency_rating"
	case "MOBO":
		sortField = "specs.tier_score"
	default:
		sortField = "specs.performance_score"
	}

	opts := options.Find().SetSort(bson.D{{Key: sortField, Value: -1}, {Key: "price_cents", Value: 1}})
	cursor, err := partRepository.collection.Find(ctx, filter, opts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var parts []*entity.Part

	for cursor.Next(ctx) {
		part := new(entity.Part)

		if err := cursor.Decode(part); err != nil {
			return nil, err
		}

		parts = append(parts, part)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return parts, nil
}
