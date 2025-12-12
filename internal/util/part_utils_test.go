package util

import (
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPartNeedToUpdate(t *testing.T) {
	tests := []struct {
		name     string
		part     *entity.Part
		expected bool
	}{
		{
			name: "part updated less than 2 hours ago should not need update",
			part: &entity.Part{
				ID:        primitive.NewObjectID(),
				Type:      entity.TypeCPU,
				Brand:     "Intel",
				Model:     "Core i7-12700K",
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
			expected: false,
		},
		{
			name: "part updated exactly 2 hours ago should need update",
			part: &entity.Part{
				ID:        primitive.NewObjectID(),
				Type:      entity.TypeCPU,
				Brand:     "Intel",
				Model:     "Core i7-12700K",
				UpdatedAt: time.Now().Add(-2 * time.Hour),
			},
			expected: true,
		},
		{
			name: "part updated more than 2 hours ago should need update",
			part: &entity.Part{
				ID:        primitive.NewObjectID(),
				Type:      entity.TypeCPU,
				Brand:     "AMD",
				Model:     "Ryzen 7 5800X",
				UpdatedAt: time.Now().Add(-3 * time.Hour),
			},
			expected: true,
		},
		{
			name: "part updated 24 hours ago should need update",
			part: &entity.Part{
				ID:        primitive.NewObjectID(),
				Type:      entity.TypeGPU,
				Brand:     "NVIDIA",
				Model:     "RTX 3080",
				UpdatedAt: time.Now().Add(-24 * time.Hour),
			},
			expected: true,
		},
		{
			name: "part updated 1 minute ago should not need update",
			part: &entity.Part{
				ID:        primitive.NewObjectID(),
				Type:      entity.TypeRAM,
				Brand:     "Corsair",
				Model:     "Vengeance 16GB",
				UpdatedAt: time.Now().Add(-1 * time.Minute),
			},
			expected: false,
		},
		{
			name:     "nil part should not need update",
			part:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PartNeedToUpdate(tt.part)
			if result != tt.expected {
				t.Errorf("PartNeedToUpdate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestPartNeedToUpdate_EdgeCases(t *testing.T) {
	t.Run("part with future timestamp should not need update", func(t *testing.T) {
		part := &entity.Part{
			ID:        primitive.NewObjectID(),
			Type:      entity.TypeCPU,
			Brand:     "Intel",
			Model:     "Core i7-12700K",
			UpdatedAt: time.Now().Add(1 * time.Hour),
		}
		result := PartNeedToUpdate(part)
		if result {
			t.Error("Part with future timestamp should not need update")
		}
	})

	t.Run("part at threshold boundary (1h 59m 59s)", func(t *testing.T) {
		part := &entity.Part{
			ID:        primitive.NewObjectID(),
			Type:      entity.TypeCPU,
			Brand:     "Intel",
			Model:     "Core i7-12700K",
			UpdatedAt: time.Now().Add(-1*time.Hour - 59*time.Minute - 59*time.Second),
		}
		result := PartNeedToUpdate(part)
		if result {
			t.Error("Part just under 2 hours should not need update")
		}
	})

	t.Run("part just over threshold (2h 0m 1s)", func(t *testing.T) {
		part := &entity.Part{
			ID:        primitive.NewObjectID(),
			Type:      entity.TypeCPU,
			Brand:     "Intel",
			Model:     "Core i7-12700K",
			UpdatedAt: time.Now().Add(-2*time.Hour - 1*time.Second),
		}
		result := PartNeedToUpdate(part)
		if !result {
			t.Error("Part just over 2 hours should need update")
		}
	})
}
