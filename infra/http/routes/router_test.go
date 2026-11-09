package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	router, err := SetupRouter(nil, nil, nil, controllers.NewHealthController(stubMongoPinger{}, true))
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

	_, err := SetupRouter(nil, nil, nil, controllers.NewHealthController(stubMongoPinger{}, true))

	if err == nil {
		t.Fatal("expected invalid TRUSTED_PROXIES to fail fast")
	}
}
