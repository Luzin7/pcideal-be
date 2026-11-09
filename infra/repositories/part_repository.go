package repositories

import (
	"context"
	"log"
	"regexp"

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

func (partRepository *PartRepository) GetAllParts(ctx context.Context, args repository.FindAllPartsArgs) ([]*entity.Part, error) {

	limit, offset := normalizePagination(args)
	opts := options.Find().
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := partRepository.collection.Find(ctx, bson.M{}, opts)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	parts := make([]*entity.Part, 0)

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

func normalizePagination(args repository.FindAllPartsArgs) (int64, int64) {
	limit := args.Limit
	if limit <= 0 {
		limit = repository.DefaultPartsPageSize
	}
	if limit > repository.MaxPartsPageSize {
		limit = repository.MaxPartsPageSize
	}

	offset := args.Offset
	if offset < 0 {
		offset = 0
	}

	return limit, offset
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

func (partRepository *PartRepository) FindPartByTypeAndBrandWithMaxPrice(ctx context.Context, args repository.FindPartByTypeAndBrandWithMaxPriceArgs) ([]*entity.Part, error) {
	opts := options.Find().SetSort(partSort(args.PartType))
	cursor, err := partRepository.collection.Find(ctx, buildPartFilter(args), opts)

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

func buildPartFilter(args repository.FindPartByTypeAndBrandWithMaxPriceArgs) bson.M {
	filter := bson.M{
		"type":        args.PartType,
		"price_cents": bson.M{"$lte": args.MaxPriceCents},
	}

	if args.Brand != "" && args.Brand != "no-preference" {
		filter["brand"] = bson.M{"$regex": regexp.QuoteMeta(args.Brand), "$options": "i"}
	}
	if args.PartType == entity.TypeCPU {
		filter["specs.socket"] = bson.M{"$exists": true, "$ne": ""}
	}
	if args.IntegratedGraphics {
		filter["specs.has_integrated_graphics"] = true
	}
	if args.MemoryType != "" {
		filter["specs.memory_type"] = args.MemoryType
	}
	if args.PartType == entity.TypePSU && args.MinPSUWatts > 0 {
		filter["specs.wattage"] = bson.M{"$gte": args.MinPSUWatts}
	}
	if args.PartType == entity.TypeMobo && args.Socket != "" {
		filter["specs.socket"] = args.Socket
	}

	return filter
}

func partSort(partType entity.PartType) bson.D {
	sortField := "specs.performance_score"

	switch partType {
	case entity.TypePSU:
		sortField = "specs.efficiency_rating"
	case entity.TypeMobo:
		sortField = "specs.tier_score"
	}

	return bson.D{{Key: sortField, Value: -1}, {Key: "price_cents", Value: 1}}
}
