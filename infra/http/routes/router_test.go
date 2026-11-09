package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Luzin7/pcideal-be/infra/http/controllers"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type stubMongoPinger struct{}

func (stubMongoPinger) Ping(context.Context, *readpref.ReadPref) error {
	return nil
}

func setupHealthRouter(t *testing.T) *gin.Engine {
	t.Helper()

	gin.SetMode(gin.TestMode)

	router, err := SetupRouter(Dependencies{
		HealthController: controllers.NewHealthController(stubMongoPinger{}, true),
	})
	if err != nil {
		t.Fatalf("expected router to be configured, got %v", err)
	}

	return router
}

func TestHealthEndpoint(t *testing.T) {
	router := setupHealthRouter(t)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", body["status"])
	}
	if body["unique_index_ready"] != true {
		t.Fatalf("expected unique_index_ready true, got %v", body["unique_index_ready"])
	}
}

func TestSetupRouterRejectsInvalidTrustedProxies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("TRUSTED_PROXIES", "not-a-cidr")

	_, err := SetupRouter(Dependencies{
		HealthController: controllers.NewHealthController(stubMongoPinger{}, true),
	})

	if err == nil {
		t.Fatal("expected invalid TRUSTED_PROXIES to fail fast")
	}
}

func TestSetupRouterRejectsInvalidAllowedOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Setenv("CORS_ALLOWED_ORIGINS", "pcideal-web.vercel.app")

	_, err := SetupRouter(Dependencies{
		HealthController: controllers.NewHealthController(stubMongoPinger{}, true),
	})

	if err == nil {
		t.Fatal("expected invalid CORS_ALLOWED_ORIGINS to fail fast")
	}
}

func TestCORSAllowsConfiguredOriginOverride(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "https://staging.example.com")
	router := setupHealthRouter(t)

	allowed := serveWithOrigin(router, "https://staging.example.com")
	if allowed.Code != http.StatusOK {
		t.Fatalf("expected configured origin to be allowed, got %d", allowed.Code)
	}
	if got := allowed.Header().Get("Access-Control-Allow-Origin"); got != "https://staging.example.com" {
		t.Fatalf("expected configured origin to be echoed, got %q", got)
	}

	rejected := serveWithOrigin(router, "https://pcideal-web.vercel.app")
	if rejected.Code != http.StatusForbidden {
		t.Fatalf("expected default origin to be rejected when overridden, got %d", rejected.Code)
	}
	if got := rejected.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allowed origin header, got %q", got)
	}
}

func TestCORSPreflightAllowsAPIKeyHeader(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	router := setupHealthRouter(t)

	recorder := servePreflight(router, "https://pcideal-web.vercel.app", "x-api-key")

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, recorder.Code)
	}
	allowHeaders := strings.ToLower(recorder.Header().Get("Access-Control-Allow-Headers"))
	if !strings.Contains(allowHeaders, "x-api-key") {
		t.Fatalf("expected x-api-key in allowed headers, got %q", allowHeaders)
	}
}

func TestCORSPreflightRejectsDisallowedOrigin(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	router := setupHealthRouter(t)

	recorder := servePreflight(router, "https://pcideal.com.br", "x-api-key")

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allowed origin header, got %q", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Headers"); got != "" {
		t.Fatalf("expected no allowed headers, got %q", got)
	}
}

func TestCORSRejectsDisallowedOrigins(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	router := setupHealthRouter(t)

	for _, origin := range []string{"http://localhost:3000", "https://pcideal.com.br", "https://example.org"} {
		recorder := serveWithOrigin(router, origin)

		if recorder.Code != http.StatusForbidden {
			t.Fatalf("expected %s to be rejected with %d, got %d", origin, http.StatusForbidden, recorder.Code)
		}
		if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("expected no allowed origin header for %s, got %q", origin, got)
		}
	}
}

func TestCORSFallsBackToDefaultOriginsWhenUnset(t *testing.T) {
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	router := setupHealthRouter(t)

	recorder := serveWithOrigin(router, "https://pcideal-web.vercel.app")

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "https://pcideal-web.vercel.app" {
		t.Fatalf("expected default origin to be echoed, got %q", got)
	}
}

func serveWithOrigin(router *gin.Engine, origin string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", origin)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}

func servePreflight(router *gin.Engine, origin string, requestedHeaders string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodOptions, "/api/builds/recommendations", nil)
	request.Header.Set("Origin", origin)
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	request.Header.Set("Access-Control-Request-Headers", requestedHeaders)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder
}
