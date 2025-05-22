package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Build struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	UserID     primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Goal       string             `bson:"goal" json:"goal"` // gaming, editing, general
	Budget     float64            `bson:"budget" json:"budget"`
	Parts      []string           `bson:"parts" json:"parts"` // Array de Part IDs
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
