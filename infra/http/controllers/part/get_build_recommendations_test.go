package part

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Luzin7/pcideal-be/infra/http/presenters"
	"github.com/Luzin7/pcideal-be/internal/domain/entity"
	"github.com/Luzin7/pcideal-be/internal/dto"
	"github.com/Luzin7/pcideal-be/internal/errors"
	"github.com/gin-gonic/gin"
)

type fakeBuildAttemptLimiter struct {
	slots    map[string]bool
	max      int
	releases int
}

func newFakeBuildAttemptLimiter(max int) *fakeBuildAttemptLimiter {
	return &fakeBuildAttemptLimiter{slots: map[string]bool{}, max: max}
}

func (f *fakeBuildAttemptLimiter) ClaimBuildAttempt(_ context.Context, attempt *entity.BuildAttempt) *errors.ErrService {
	window := attempt.CreatedAt.UTC().Truncate(time.Hour)

	for slot := 0; slot < f.max; slot++ {
		key := buildAttemptSlotKey(attempt.IP, window, slot)
		if f.slots[key] {
			continue
		}
		f.slots[key] = true
		attempt.WindowStart = window
		attempt.Slot = slot
		return nil
	}

	return errors.ErrBuildAttemptLimit()
}

func (f *fakeBuildAttemptLimiter) ReleaseBuildAttempt(ctx context.Context, attempt *entity.BuildAttempt) *errors.ErrService {
	if ctx.Err() != nil {
		return errors.New("context canceled", http.StatusInternalServerError)
	}

	delete(f.slots, buildAttemptSlotKey(attempt.IP, attempt.WindowStart, attempt.Slot))
	f.releases++
	return nil
}

func buildAttemptSlotKey(ip string, window time.Time, slot int) string {
	return ip + "|" + window.Format(time.RFC3339) + "|" + strconv.Itoa(slot)
}

type fakeBuildRecommendationsGenerator struct {
	calls int
	err   *errors.ErrService
}

func (f *fakeBuildRecommendationsGenerator) Execute(context.Context, dto.GenerateBuildRecommendationsDTO) (*presenters.RecommendedBuildsPresenter, *errors.ErrService) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}

	return &presenters.RecommendedBuildsPresenter{}, nil
}

func performRecommendationRequest(t *testing.T, controller *GetBuildRecommendationsController, ctx context.Context) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/builds/recommendations", controller.Handle)

	request := httptest.NewRequest(
		http.MethodPost,
		"/builds/recommendations",
		strings.NewReader(`{"usage_type":"GAMING","budget_cents":500000}`),
	).WithContext(ctx)
	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func TestRecommendationsDoNotGenerateAfterLimitReached(t *testing.T) {
	limiter := newFakeBuildAttemptLimiter(4)
	generator := &fakeBuildRecommendationsGenerator{}
	controller := NewGetBuildRecommendationsController(generator, limiter)

	for i := 0; i < 4; i++ {
		if recorder := performRecommendationRequest(t, controller, context.Background()); recorder.Code != http.StatusOK {
			t.Fatalf("request %d expected 200, got %d", i+1, recorder.Code)
		}
	}

	recorder := performRecommendationRequest(t, controller, context.Background())

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", recorder.Code)
	}
	if generator.calls != 4 {
		t.Fatalf("expected generator called 4 times, got %d", generator.calls)
	}
}

func TestGenerationFailureReleasesClaimedSlot(t *testing.T) {
	limiter := newFakeBuildAttemptLimiter(1)
	generator := &fakeBuildRecommendationsGenerator{err: errors.ErrNoCompatibleParts("CPU")}
	controller := NewGetBuildRecommendationsController(generator, limiter)

	recorder := performRecommendationRequest(t, controller, context.Background())
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
	if limiter.releases != 1 {
		t.Fatalf("expected 1 release, got %d", limiter.releases)
	}

	generator.err = nil

	recorder = performRecommendationRequest(t, controller, context.Background())
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected retry to succeed, got %d", recorder.Code)
	}
}

func TestCanceledRequestContextStillReleasesClaimedSlot(t *testing.T) {
	limiter := newFakeBuildAttemptLimiter(1)
	generator := &fakeBuildRecommendationsGenerator{err: errors.ErrNoCompatibleParts("CPU")}
	controller := NewGetBuildRecommendationsController(generator, limiter)

	requestContext, cancelRequest := context.WithCancel(context.Background())
	cancelRequest()

	recorder := performRecommendationRequest(t, controller, requestContext)
	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", recorder.Code)
	}
	if len(limiter.slots) != 0 {
		t.Fatalf("expected slot released, %d still held", len(limiter.slots))
	}

	generator.err = nil

	recorder = performRecommendationRequest(t, controller, context.Background())
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected retry to succeed, got %d", recorder.Code)
	}
}

func TestSuccessfulClaimPersists(t *testing.T) {
	limiter := newFakeBuildAttemptLimiter(4)
	generator := &fakeBuildRecommendationsGenerator{}
	controller := NewGetBuildRecommendationsController(generator, limiter)

	recorder := performRecommendationRequest(t, controller, context.Background())

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if len(limiter.slots) != 1 {
		t.Fatalf("expected one persisted claim, got %d", len(limiter.slots))
	}
	if limiter.releases != 0 {
		t.Fatalf("expected no release, got %d", limiter.releases)
	}
}
