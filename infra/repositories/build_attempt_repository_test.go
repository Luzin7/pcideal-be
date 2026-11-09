package repositories

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func duplicateKeyError() error {
	return mongo.WriteException{WriteErrors: mongo.WriteErrors{{Code: 11000}}}
}

func TestClaimAttemptSlotRetriesOnDuplicate(t *testing.T) {
	attempt := &entity.BuildAttempt{CreatedAt: time.Date(2026, 11, 9, 14, 35, 0, 0, time.UTC)}
	calls := 0

	err := claimAttemptSlot(context.Background(), attempt, func(_ context.Context, a *entity.BuildAttempt) error {
		calls++
		if a.Slot < 1 {
			return duplicateKeyError()
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 insert attempts, got %d", calls)
	}
	if attempt.Slot != 1 {
		t.Fatalf("expected slot 1, got %d", attempt.Slot)
	}
	expectedWindow := time.Date(2026, 11, 9, 14, 0, 0, 0, time.UTC)
	if !attempt.WindowStart.Equal(expectedWindow) {
		t.Fatalf("expected window start %v, got %v", expectedWindow, attempt.WindowStart)
	}
}

func TestClaimAttemptSlotLimitReached(t *testing.T) {
	attempt := &entity.BuildAttempt{CreatedAt: time.Now()}

	err := claimAttemptSlot(context.Background(), attempt, func(context.Context, *entity.BuildAttempt) error {
		return duplicateKeyError()
	})

	if !errors.Is(err, repository.ErrBuildAttemptLimitReached) {
		t.Fatalf("expected ErrBuildAttemptLimitReached, got %v", err)
	}
}

func TestClaimAttemptSlotPropagatesOtherErrors(t *testing.T) {
	attempt := &entity.BuildAttempt{CreatedAt: time.Now()}
	insertErr := errors.New("network down")

	err := claimAttemptSlot(context.Background(), attempt, func(context.Context, *entity.BuildAttempt) error {
		return insertErr
	})

	if !errors.Is(err, insertErr) {
		t.Fatalf("expected %v, got %v", insertErr, err)
	}
}

func TestClaimAttemptSlotConcurrentDuplicateOnlyOneWins(t *testing.T) {
	var mutex sync.Mutex
	taken := map[string]bool{}
	insert := func(_ context.Context, attempt *entity.BuildAttempt) error {
		mutex.Lock()
		defer mutex.Unlock()

		key := fmt.Sprintf("%s|%d", attempt.WindowStart, attempt.Slot)
		if taken[key] {
			return duplicateKeyError()
		}
		taken[key] = true
		return nil
	}

	window := time.Date(2026, 11, 9, 14, 35, 0, 0, time.UTC)
	for slot := 1; slot < maxBuildAttemptsPerWindow; slot++ {
		taken[fmt.Sprintf("%s|%d", window.Truncate(time.Hour), slot)] = true
	}

	results := make(chan error, 2)
	start := make(chan struct{})

	for i := 0; i < 2; i++ {
		go func() {
			<-start
			results <- claimAttemptSlot(context.Background(), &entity.BuildAttempt{CreatedAt: window}, insert)
		}()
	}
	close(start)

	var successes, limits int
	for i := 0; i < 2; i++ {
		err := <-results
		if err == nil {
			successes++
			continue
		}
		if errors.Is(err, repository.ErrBuildAttemptLimitReached) {
			limits++
		}
	}

	if successes != 1 || limits != 1 {
		t.Fatalf("expected one success and one limit, got %d success and %d limit", successes, limits)
	}
}

func TestBuildAttemptReleaseFilterTargetsClaimedSlot(t *testing.T) {
	window := time.Date(2026, 11, 9, 14, 0, 0, 0, time.UTC)
	attempt := &entity.BuildAttempt{IP: "203.0.113.7", WindowStart: window, Slot: 2}

	filter := buildAttemptReleaseFilter(attempt)

	expected := bson.M{"ip": "203.0.113.7", "window_start": window, "slot": 2}
	if fmt.Sprint(filter) != fmt.Sprint(expected) {
		t.Fatalf("expected filter %v, got %v", expected, filter)
	}
}
