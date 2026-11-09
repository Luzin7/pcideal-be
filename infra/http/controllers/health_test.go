package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type stubMongoPinger struct {
	err error
}

func (s stubMongoPinger) Ping(context.Context, *readpref.ReadPref) error {
	return s.err
}

func performHealthRequest(t *testing.T, controller *HealthController) *httptest.ResponseRecorder {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/health", controller.Handle)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	return recorder
}

func TestHealthReturnsOKWhenMongoIsReachable(t *testing.T) {
	recorder := performHealthRequest(t, NewHealthController(stubMongoPinger{}, true))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
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

func TestHealthReturnsServiceUnavailableWhenMongoIsDown(t *testing.T) {
	recorder := performHealthRequest(t, NewHealthController(stubMongoPinger{err: errors.New("mongo down")}, false))

	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", recorder.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("expected JSON body, got %v", err)
	}
	if body["status"] != "unavailable" {
		t.Fatalf("expected status unavailable, got %v", body["status"])
	}
	if body["unique_index_ready"] != false {
		t.Fatalf("expected unique_index_ready false, got %v", body["unique_index_ready"])
	}
}

func TestHealthReportsUniqueIndexWithoutFailingReadiness(t *testing.T) {
	recorder := performHealthRequest(t, NewHealthController(stubMongoPinger{}, false))

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when only the unique index is missing, got %d", recorder.Code)
	}
}
