package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Price struct {
	ID          string    `bson:"_id" json:"id"`
	PartID      string    `bson:"part_id" json:"part_id"`
	Price       int64     `bson:"price_cents" json:"price_cents"`
	Store       string    `bson:"store" json:"store"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}

func NewPrice(partID string, price int64, store string) *Price {
	return &Price{
		ID:          primitive.NewObjectID().Hex(),
		PartID:      partID,
		Price:       price,
		Store:       store,
		LastUpdated: time.Now(),
	}
}
