package part

import (
	"testing"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestChooseBestPartBreaksTieByID(t *testing.T) {
	lowID, err := primitive.ObjectIDFromHex("000000000000000000000001")
	if err != nil {
		t.Fatalf("failed to build low id: %v", err)
	}
	highID, err := primitive.ObjectIDFromHex("ffffffffffffffffffffffff")
	if err != nil {
		t.Fatalf("failed to build high id: %v", err)
	}

	high := &entity.Part{ID: highID, Specs: entity.Specs{PerformanceScore: 80}, PriceCents: 50000}
	low := &entity.Part{ID: lowID, Specs: entity.Specs{PerformanceScore: 80}, PriceCents: 50000}

	best := chooseBestPart([]*entity.Part{high, low}, partRanking{score: performanceScore})

	if best == nil || best.ID != lowID {
		t.Fatalf("expected low id to win, got %v", best)
	}
}
