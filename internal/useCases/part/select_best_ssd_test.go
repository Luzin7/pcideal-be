package part

import (
	"context"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSelectBestSSDUsesPerformanceScore(t *testing.T) {
	uc := NewSelectBestSSDUseCase(newMockUpdatePartsUseCase())

	fast := &entity.Part{
		ID:         primitive.NewObjectID(),
		Type:       entity.TypeSSD,
		Model:      "Fast",
		Specs:      entity.Specs{PerformanceScore: 90, EfficiencyRating: 1},
		PriceCents: 50000,
		UpdatedAt:  time.Now(),
	}
	rated := &entity.Part{
		ID:         primitive.NewObjectID(),
		Type:       entity.TypeSSD,
		Model:      "Rated",
		Specs:      entity.Specs{PerformanceScore: 40, EfficiencyRating: 5},
		PriceCents: 50000,
		UpdatedAt:  time.Now(),
	}

	result, err := uc.Execute(context.Background(), SelectBestSSDArgs{
		ssds: []*entity.Part{rated, fast},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Model != fast.Model {
		t.Fatalf("expected %s, got %s", fast.Model, result.Model)
	}
}

func TestSelectBestSSDBreaksTieByCapacity(t *testing.T) {
	uc := NewSelectBestSSDUseCase(newMockUpdatePartsUseCase())

	small := &entity.Part{
		ID:         primitive.NewObjectID(),
		Type:       entity.TypeSSD,
		Model:      "Small",
		Specs:      entity.Specs{PerformanceScore: 80, CapacityGB: 480},
		PriceCents: 50000,
		UpdatedAt:  time.Now(),
	}
	large := &entity.Part{
		ID:         primitive.NewObjectID(),
		Type:       entity.TypeSSD,
		Model:      "Large",
		Specs:      entity.Specs{PerformanceScore: 80, CapacityGB: 960},
		PriceCents: 50000,
		UpdatedAt:  time.Now(),
	}

	result, err := uc.Execute(context.Background(), SelectBestSSDArgs{
		ssds: []*entity.Part{small, large},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Model != large.Model {
		t.Fatalf("expected %s, got %s", large.Model, result.Model)
	}
}
