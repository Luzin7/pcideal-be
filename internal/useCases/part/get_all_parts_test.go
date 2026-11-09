package part

import (
	"context"
	"testing"

	"github.com/Luzin7/pcideal-be/internal/domain/repository"
)

func TestGetAllPartsReturnsEmptySlice(t *testing.T) {
	uc := NewGetAllPartsUseCase(emptyPartRepository{})

	parts, err := uc.Execute(context.Background(), repository.FindAllPartsArgs{Limit: 50})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if parts == nil {
		t.Fatal("expected empty slice, got nil")
	}
	if len(parts) != 0 {
		t.Fatalf("expected no parts, got %d", len(parts))
	}
}
