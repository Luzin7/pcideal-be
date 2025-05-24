package models

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewPart(t *testing.T) {
	specs := specs{
		Socket:     "AM4",
		TDP:        65,
		ClockSpeed: 3.6,
	}

	part := NewPart("CPU", "AMD", "Ryzen 5 5600X", specs, 120000, "https://kabum.com.br/ryzen5600x", "Kabum")

	if part == nil {
		t.Fatal("expected part to not be nil")
	}

	if part.Type != "CPU" {
		t.Errorf("expected Type to be 'CPU', got '%s'", part.Type)
	}
	if part.Brand != "AMD" {
		t.Errorf("expected Brand to be 'AMD', got '%s'", part.Brand)
	}
	if part.Model != "Ryzen 5 5600X" {
		t.Errorf("expected Model to be 'Ryzen 5 5600X', got '%s'", part.Model)
	}
	if part.PriceCents != 120000 {
		t.Errorf("expected PriceCents to be 120000, got %d", part.PriceCents)
	}
	if part.URL != "https://kabum.com.br/ryzen5600x" {
		t.Errorf("expected URL to be correct, got %s", part.URL)
	}
	if part.Store != "Kabum" {
		t.Errorf("expected Store to be Kabum, got %s", part.Store)
	}

	if part.ID == primitive.NilObjectID {
		t.Errorf("expected non-empty ID")
	}
}
