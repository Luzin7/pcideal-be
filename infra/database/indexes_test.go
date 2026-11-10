package database

import (
	"errors"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestPartIndexesCoverQueryShape(t *testing.T) {
	var performanceIndex *mongo.IndexModel
	var storeURLIndex *mongo.IndexModel

	indexes := partIndexes()
	for i := range indexes {
		index := &indexes[i]
		if index.Options.Name == nil {
			continue
		}

		switch *index.Options.Name {
		case "parts_type_price_performance":
			performanceIndex = index
		case "parts_store_url_unique":
			storeURLIndex = index
		}
	}

	if performanceIndex == nil {
		t.Fatal("expected parts_type_price_performance index")
	}
	assertIndexKeys(t, performanceIndex, []string{"type", "price_cents", "specs.performance_score"})

	if storeURLIndex == nil {
		t.Fatal("expected parts_store_url_unique index")
	}
	assertIndexKeys(t, storeURLIndex, []string{"store", "url"})

	if storeURLIndex.Options.Unique == nil || !*storeURLIndex.Options.Unique {
		t.Fatal("expected store/url index to be unique")
	}
}

func TestBuildAttemptIndexIsUnique(t *testing.T) {
	indexes := buildAttemptIndexes()
	if len(indexes) != 1 {
		t.Fatalf("expected 1 build_attempts index, got %d", len(indexes))
	}

	assertIndexKeys(t, &indexes[0], []string{"ip", "window_start", "slot"})

	if indexes[0].Options.Unique == nil || !*indexes[0].Options.Unique {
		t.Fatal("expected build_attempts index to be unique")
	}
}

func TestHasUniqueStoreURLIndex(t *testing.T) {
	storeURLKey := bson.D{{Key: "store", Value: 1}, {Key: "url", Value: 1}}

	tests := []struct {
		name    string
		indexes []IndexSpec
		want    bool
	}{
		{
			name:    "unique index present",
			indexes: []IndexSpec{{Unique: true, Key: storeURLKey}},
			want:    true,
		},
		{
			name:    "unique index with reversed key order",
			indexes: []IndexSpec{{Unique: true, Key: bson.D{{Key: "url", Value: 1}, {Key: "store", Value: 1}}}},
			want:    true,
		},
		{
			name:    "unique index on other keys",
			indexes: []IndexSpec{{Unique: true, Key: bson.D{{Key: "url", Value: 1}}}},
			want:    false,
		},
		{
			name:    "index missing",
			indexes: []IndexSpec{},
			want:    false,
		},
		{
			name:    "index present but not unique",
			indexes: []IndexSpec{{Unique: false, Key: storeURLKey}},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasUniqueStoreURLIndex(tt.indexes); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestRequireUniqueStoreURLIndex(t *testing.T) {
	checkErr := errors.New("list parts indexes: boom")

	tests := []struct {
		name        string
		ready       bool
		checkErr    error
		wantErr     bool
		wantMessage string
	}{
		{
			name:  "ready",
			ready: true,
		},
		{
			name:        "not ready",
			ready:       false,
			wantErr:     true,
			wantMessage: StoreURLDuplicateAuditCommand,
		},
		{
			name:     "check failed",
			checkErr: checkErr,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RequireUniqueStoreURLIndex(tt.ready, tt.checkErr)
			if tt.wantErr != (err != nil) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			if !tt.wantErr {
				return
			}

			var target *UniqueIndexNotReadyError
			if !errors.As(err, &target) {
				t.Fatalf("expected *UniqueIndexNotReadyError, got %T", err)
			}

			if tt.checkErr != nil && !errors.Is(err, tt.checkErr) {
				t.Fatalf("expected cause %v, got %v", tt.checkErr, err)
			}

			if tt.wantMessage != "" && !strings.Contains(err.Error(), tt.wantMessage) {
				t.Fatalf("expected message to contain %q, got %q", tt.wantMessage, err.Error())
			}
		})
	}
}

func assertIndexKeys(t *testing.T, index *mongo.IndexModel, expected []string) {
	t.Helper()

	keys, ok := index.Keys.(bson.D)
	if !ok {
		t.Fatalf("expected bson.D keys, got %T", index.Keys)
	}

	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}

	for i, key := range expected {
		if keys[i].Key != key {
			t.Fatalf("expected key %q at position %d, got %q", key, i, keys[i].Key)
		}
	}
}
