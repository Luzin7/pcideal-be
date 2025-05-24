package models

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewPrice(t *testing.T) {
	price := NewPrice("1", 10000, "TestStore")

	if price == nil {
		t.Fatal("expected price to not be nil")
	}

	if price.ID == primitive.NilObjectID {
		t.Errorf("expected price ID to be set")
	}

	if price.PartID != "1" {
		t.Errorf("expected PartID to be set")
	}

	if price.Price != 10000 {
		t.Errorf("expected price to be 10000, got %d", price.Price)
	}

	if price.Store != "TestStore" {
		t.Errorf("expected store to be 'TestStore', got '%s'", price.Store)
	}
}
