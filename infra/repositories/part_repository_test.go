package repositories

import (
	"testing"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func TestBuildPartFilterMoboUsesSocket(t *testing.T) {
	filter := buildPartFilter(repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypeMobo,
		Socket:        "AM4",
		MaxPriceCents: 100000,
	})

	if filter["type"] != entity.TypeMobo {
		t.Fatalf("expected type %v, got %v", entity.TypeMobo, filter["type"])
	}
	if filter["specs.socket"] != "AM4" {
		t.Fatalf("expected socket AM4, got %v", filter["specs.socket"])
	}
}

func TestPartSortMoboUsesTierScore(t *testing.T) {
	sort := partSort(entity.TypeMobo)

	if len(sort) == 0 {
		t.Fatal("expected sort keys")
	}
	if sort[0].Key != "specs.tier_score" || sort[0].Value != -1 {
		t.Fatalf("expected specs.tier_score desc, got %v", sort[0])
	}
}

func TestNormalizePaginationBounds(t *testing.T) {
	tests := []struct {
		name       string
		args       repository.FindAllPartsArgs
		wantLimit  int64
		wantOffset int64
	}{
		{name: "defaults zero limit", args: repository.FindAllPartsArgs{}, wantLimit: repository.DefaultPartsPageSize, wantOffset: 0},
		{name: "caps oversized limit", args: repository.FindAllPartsArgs{Limit: repository.MaxPartsPageSize + 1}, wantLimit: repository.MaxPartsPageSize, wantOffset: 0},
		{name: "keeps valid values", args: repository.FindAllPartsArgs{Limit: 10, Offset: 25}, wantLimit: 10, wantOffset: 25},
		{name: "floors negative offset", args: repository.FindAllPartsArgs{Limit: 10, Offset: -1}, wantLimit: 10, wantOffset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit, offset := normalizePagination(tt.args)

			if limit != tt.wantLimit || offset != tt.wantOffset {
				t.Fatalf("expected limit %d offset %d, got %d/%d", tt.wantLimit, tt.wantOffset, limit, offset)
			}
		})
	}
}

func TestBuildPartFilterEscapesBrandRegex(t *testing.T) {
	filter := buildPartFilter(repository.FindPartByTypeAndBrandWithMaxPriceArgs{
		PartType:      entity.TypeGPU,
		Brand:         "ASUS (ROG)",
		MaxPriceCents: 100000,
	})

	brandFilter, ok := filter["brand"].(bson.M)
	if !ok {
		t.Fatalf("expected bson.M brand filter, got %T", filter["brand"])
	}
	if brandFilter["$regex"] != `ASUS \(ROG\)` {
		t.Fatalf("expected escaped regex, got %v", brandFilter["$regex"])
	}
	if brandFilter["$options"] != "i" {
		t.Fatalf("expected case-insensitive option, got %v", brandFilter["$options"])
	}
}
