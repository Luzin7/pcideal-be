package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Price struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	PartID    string             `bson:"part_id" json:"part_id"`
	Price     int64              `bson:"price_cents" json:"price_cents"`
	Store     string             `bson:"store" json:"store"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewPrice(partID string, price int64, store string) *Price {
	return &Price{
		ID:        primitive.NewObjectID(),
		PartID:    partID,
		Price:     price,
		Store:     store,
		UpdatedAt: time.Now(),
	}
}
