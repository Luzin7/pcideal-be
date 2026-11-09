package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const StoreURLIndexName = "parts_store_url_unique"

const StoreURLDuplicateAuditCommand = `db.parts.aggregate([{ $group: { _id: { store: "$store", url: "$url" }, count: { $sum: 1 } } }, { $match: { count: { $gt: 1 } } }])`

type IndexSpec struct {
	Unique bool   `bson:"unique"`
	Key    bson.D `bson:"key"`
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	if _, err := db.Collection("parts").Indexes().CreateMany(ctx, nonUniquePartIndexes()); err != nil {
		return fmt.Errorf("ensure parts indexes: %w", err)
	}

	if _, err := db.Collection("build_attempts").Indexes().CreateMany(ctx, buildAttemptIndexes()); err != nil {
		return fmt.Errorf("ensure build_attempts indexes: %w", err)
	}

	if _, err := db.Collection("parts").Indexes().CreateOne(ctx, storeURLUniqueIndex()); err != nil {
		log.Printf(
			"[database] non-fatal: failed to create unique index %s: %v; run duplicate audit: %s",
			StoreURLIndexName, err, StoreURLDuplicateAuditCommand,
		)
	}

	return nil
}

func UniqueStoreURLIndexReady(ctx context.Context, db *mongo.Database) (bool, error) {
	cursor, err := db.Collection("parts").Indexes().List(ctx)
	if err != nil {
		return false, fmt.Errorf("list parts indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexes []IndexSpec
	if err := cursor.All(ctx, &indexes); err != nil {
		return false, fmt.Errorf("decode parts indexes: %w", err)
	}

	return hasUniqueStoreURLIndex(indexes), nil
}

func hasUniqueStoreURLIndex(indexes []IndexSpec) bool {
	for _, index := range indexes {
		if index.Unique && isStoreURLIndexKey(index.Key) {
			return true
		}
	}

	return false
}

func isStoreURLIndexKey(key bson.D) bool {
	if len(key) != 2 {
		return false
	}

	return (key[0].Key == "store" && key[1].Key == "url") ||
		(key[0].Key == "url" && key[1].Key == "store")
}

func partIndexes() []mongo.IndexModel {
	return append(nonUniquePartIndexes(), storeURLUniqueIndex())
}

func nonUniquePartIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "price_cents", Value: 1},
				{Key: "specs.performance_score", Value: -1},
			},
			Options: options.Index().SetName("parts_type_price_performance"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "price_cents", Value: 1},
				{Key: "specs.tier_score", Value: -1},
			},
			Options: options.Index().SetName("parts_type_price_tier"),
		},
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "price_cents", Value: 1},
				{Key: "specs.efficiency_rating", Value: -1},
			},
			Options: options.Index().SetName("parts_type_price_efficiency"),
		},
	}
}

func storeURLUniqueIndex() mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.D{
			{Key: "store", Value: 1},
			{Key: "url", Value: 1},
		},
		Options: options.Index().SetName(StoreURLIndexName).SetUnique(true),
	}
}

func buildAttemptIndexes() []mongo.IndexModel {
	return []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "ip", Value: 1},
				{Key: "window_start", Value: 1},
				{Key: "slot", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetPartialFilterExpression(bson.M{"window_start": bson.M{"$exists": true}}),
		},
	}
}
