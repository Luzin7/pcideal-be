package main

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPickCanonical(t *testing.T) {
	olderTimestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	newerTimestamp := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		documents []entity.Part
		wantID    string
	}{
		{
			name: "newest updated_at wins",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(newerTimestamp)),
			},
			wantID: "000000000000000000000002",
		},
		{
			name: "newest updated_at wins despite older id",
			documents: []entity.Part{
				newTestPart(withID("00000000000000000000000f"), withUpdatedAt(newerTimestamp)),
				newTestPart(withID("000000000000000000000010"), withUpdatedAt(olderTimestamp)),
			},
			wantID: "00000000000000000000000f",
		},
		{
			name: "missing updated_at loses to dated document",
			documents: []entity.Part{
				newTestPart(withID("00000000000000000000000f"), withoutUpdatedAt()),
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
			},
			wantID: "000000000000000000000001",
		},
		{
			name: "tie uses newest id",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(olderTimestamp)),
			},
			wantID: "000000000000000000000002",
		},
		{
			name: "both missing updated_at uses newest id",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000010"), withoutUpdatedAt()),
				newTestPart(withID("00000000000000000000000f"), withoutUpdatedAt()),
			},
			wantID: "000000000000000000000010",
		},
		{
			name:      "single document is canonical",
			documents: []entity.Part{newTestPart(withID("000000000000000000000003"))},
			wantID:    "000000000000000000000003",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickCanonical(tt.documents)
			if got.ID.Hex() != tt.wantID {
				t.Fatalf("expected canonical %s, got %s", tt.wantID, got.ID.Hex())
			}
		})
	}
}

func TestPlanDeletions(t *testing.T) {
	olderTimestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	newerTimestamp := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		documents []entity.Part
		wantIDs   []string
	}{
		{
			name: "deletes document with older updated_at",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(newerTimestamp)),
			},
			wantIDs: []string{"000000000000000000000001"},
		},
		{
			name: "three documents delete two",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000003"), withUpdatedAt(newerTimestamp)),
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(olderTimestamp)),
			},
			wantIDs: []string{"000000000000000000000001", "000000000000000000000002"},
		},
		{
			name: "tie deletes older id",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(olderTimestamp)),
			},
			wantIDs: []string{"000000000000000000000001"},
		},
		{
			name: "deletes document missing updated_at",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withoutUpdatedAt()),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(olderTimestamp)),
			},
			wantIDs: []string{"000000000000000000000001"},
		},
		{
			name:      "single document deletes nothing",
			documents: []entity.Part{newTestPart(withID("000000000000000000000001"))},
			wantIDs:   nil,
		},
		{
			name:      "empty group deletes nothing",
			documents: nil,
			wantIDs:   nil,
		},
		{
			name: "identical documents keep newest id",
			documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000003"), withUpdatedAt(olderTimestamp)),
			},
			wantIDs: []string{"000000000000000000000001", "000000000000000000000002"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := planDeletions(tt.documents)
			if len(got) != len(tt.wantIDs) {
				t.Fatalf("expected %d deletions, got %d", len(tt.wantIDs), len(got))
			}

			for i, wantID := range tt.wantIDs {
				if got[i].Hex() != wantID {
					t.Fatalf("expected deletion %s at position %d, got %s", wantID, i, got[i].Hex())
				}
			}
		})
	}
}

func TestDuplicateDeleteFilter(t *testing.T) {
	group := duplicateGroup{
		Store: "kabum",
		URL:   "https://www.kabum.com.br/produto/1",
	}
	ids := []primitive.ObjectID{
		objectID("000000000000000000000001"),
		objectID("000000000000000000000002"),
	}

	filter := duplicateDeleteFilter(group, ids)

	if filter["store"] != group.Store {
		t.Fatalf("expected store %q, got %v", group.Store, filter["store"])
	}
	if filter["url"] != group.URL {
		t.Fatalf("expected url %q, got %v", group.URL, filter["url"])
	}

	idFilter, ok := filter["_id"].(bson.M)
	if !ok {
		t.Fatalf("expected bson.M id filter, got %T", filter["_id"])
	}

	inIDs, ok := idFilter["$in"].([]primitive.ObjectID)
	if !ok {
		t.Fatalf("expected []primitive.ObjectID in filter, got %T", idFilter["$in"])
	}
	if len(inIDs) != len(ids) {
		t.Fatalf("expected %d ids, got %d", len(ids), len(inIDs))
	}
	if inIDs[0].Hex() != "000000000000000000000001" || inIDs[1].Hex() != "000000000000000000000002" {
		t.Fatalf("expected ids %v, got %v", ids, inIDs)
	}
}

func TestIsIndexConflict(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "index options conflict",
			err:  fmt.Errorf("create unique index: %w", mongo.CommandError{Code: 85}),
			want: true,
		},
		{
			name: "index key specs conflict",
			err:  fmt.Errorf("create unique index: %w", mongo.CommandError{Code: 86}),
			want: true,
		},
		{
			name: "duplicate key is not an index conflict",
			err:  fmt.Errorf("insert: %w", mongo.CommandError{Code: 11000}),
			want: false,
		},
		{
			name: "unrelated error",
			err:  errors.New("boom"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isIndexConflict(tt.err); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestAuditedDocumentPushUsesSharedFields(t *testing.T) {
	projection := partProjection()
	push := auditedDocumentPush()

	if len(push) != len(partFields)+1 {
		t.Fatalf("expected %d push fields, got %d", len(partFields)+1, len(push))
	}
	if push[0].Key != "_id" || push[0].Value != "$_id" {
		t.Fatalf("expected _id at position 0, got %v", push[0])
	}

	for index, field := range partFields {
		element := push[index+1]
		if element.Key != field || element.Value != "$"+field {
			t.Fatalf("expected push field %s, got %v", field, element)
		}
		if projection[field] != 1 {
			t.Fatalf("expected projection for %s, got %v", field, projection[field])
		}
	}
}

func TestHasDifferenceSingleDocument(t *testing.T) {
	documents := []entity.Part{newTestPart(withBrand("AMD"))}

	if hasDifference(documents, func(part entity.Part) string { return part.Brand }) {
		t.Fatal("expected no difference for a single document")
	}
}

func TestCanonicalIDs(t *testing.T) {
	olderTimestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	newerTimestamp := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)

	groups := []duplicateGroup{
		{
			Store: "kabum",
			URL:   "https://www.kabum.com.br/produto/1",
			Documents: []entity.Part{
				newTestPart(withID("000000000000000000000001"), withUpdatedAt(olderTimestamp)),
				newTestPart(withID("000000000000000000000002"), withUpdatedAt(newerTimestamp)),
			},
		},
		{
			Store: "kabum",
			URL:   "https://www.kabum.com.br/produto/2",
			Documents: []entity.Part{
				newTestPart(withID("000000000000000000000003"), withoutUpdatedAt()),
				newTestPart(withID("000000000000000000000004"), withUpdatedAt(olderTimestamp)),
			},
		},
	}

	got := canonicalIDs(groups)

	if len(got) != 2 {
		t.Fatalf("expected 2 canonical ids, got %d", len(got))
	}
	if got[0].Hex() != "000000000000000000000002" {
		t.Fatalf("expected canonical 000000000000000000000002, got %s", got[0].Hex())
	}
	if got[1].Hex() != "000000000000000000000004" {
		t.Fatalf("expected canonical 000000000000000000000004, got %s", got[1].Hex())
	}
}

func TestCanonicalRetentionComplete(t *testing.T) {
	tests := []struct {
		name     string
		expected int
		retained int64
		want     bool
	}{
		{name: "all canonicals retained", expected: 2, retained: 2, want: true},
		{name: "one canonical missing", expected: 2, retained: 1, want: false},
		{name: "no groups", expected: 0, retained: 0, want: true},
		{name: "unexpected canonical", expected: 0, retained: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canonicalRetentionComplete(tt.expected, tt.retained); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestDifferingFields(t *testing.T) {
	baseTimestamp := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		documents []entity.Part
		want      []string
	}{
		{
			name: "identical documents differ in nothing",
			documents: []entity.Part{
				newTestPart(),
				newTestPart(),
			},
			want: nil,
		},
		{
			name: "price differs",
			documents: []entity.Part{
				newTestPart(),
				newTestPart(withPriceCents(99900)),
			},
			want: []string{"price_cents"},
		},
		{
			name: "brand and model differ",
			documents: []entity.Part{
				newTestPart(),
				newTestPart(withBrand("Intel"), withModel("Core i5-12400")),
			},
			want: []string{"brand", "model"},
		},
		{
			name: "updated_at differs",
			documents: []entity.Part{
				newTestPart(),
				newTestPart(withUpdatedAt(baseTimestamp.Add(time.Hour))),
			},
			want: []string{"updated_at"},
		},
		{
			name: "all fields differ in fixed order",
			documents: []entity.Part{
				newTestPart(),
				newTestPart(
					withPriceCents(99900),
					withBrand("Intel"),
					withModel("Core i5-12400"),
					withUpdatedAt(baseTimestamp.Add(time.Hour)),
				),
			},
			want: []string{"price_cents", "brand", "model", "updated_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := differingFields(tt.documents)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Fatalf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func TestCheckFailed(t *testing.T) {
	tests := []struct {
		name    string
		summary auditSummary
		want    bool
	}{
		{name: "ready to apply", summary: auditSummary{}, want: false},
		{name: "duplicate groups present", summary: auditSummary{Groups: 1}, want: true},
		{name: "missing store or url present", summary: auditSummary{MissingKeys: 1}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkFailed(tt.summary); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestRunRejectsInvalidApplyFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "apply without confirmation", args: []string{"--apply"}},
		{name: "apply with wrong confirmation", args: []string{"--apply", "--confirm=BAD"}},
		{name: "unknown flag", args: []string{"--unknown"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.args, io.Discard, io.Discard); got != 2 {
				t.Fatalf("expected exit 2, got %d", got)
			}
		})
	}
}

func TestValidateApplyMode(t *testing.T) {
	tests := []struct {
		name         string
		apply        bool
		confirmation string
		wantErr      bool
	}{
		{name: "dry run needs no confirmation", apply: false, confirmation: "", wantErr: false},
		{name: "apply with confirmation is valid", apply: true, confirmation: deleteConfirmation, wantErr: false},
		{name: "apply without confirmation fails", apply: true, confirmation: "", wantErr: true},
		{name: "apply with wrong confirmation fails", apply: true, confirmation: "DELETE", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateApplyMode(tt.apply, tt.confirmation)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				return
			}

			var confirmationErr *MissingConfirmationError
			if !errors.As(err, &confirmationErr) {
				t.Fatalf("expected MissingConfirmationError, got %T", err)
			}
		})
	}
}

type partOption func(*entity.Part)

func newTestPart(options ...partOption) entity.Part {
	part := entity.Part{
		ID:         objectID("000000000000000000000001"),
		Store:      "kabum",
		URL:        "https://www.kabum.com.br/produto/1",
		PriceCents: 100000,
		Brand:      "AMD",
		Model:      "Ryzen 5 5600",
		UpdatedAt:  time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
	}

	for _, option := range options {
		option(&part)
	}

	return part
}

func withID(hex string) partOption {
	return func(part *entity.Part) {
		part.ID = objectID(hex)
	}
}

func withUpdatedAt(updatedAt time.Time) partOption {
	return func(part *entity.Part) {
		part.UpdatedAt = updatedAt
	}
}

func withoutUpdatedAt() partOption {
	return func(part *entity.Part) {
		part.UpdatedAt = time.Time{}
	}
}

func withPriceCents(priceCents int64) partOption {
	return func(part *entity.Part) {
		part.PriceCents = priceCents
	}
}

func withBrand(brand string) partOption {
	return func(part *entity.Part) {
		part.Brand = brand
	}
}

func withModel(model string) partOption {
	return func(part *entity.Part) {
		part.Model = model
	}
}

func objectID(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}

	return id
}
