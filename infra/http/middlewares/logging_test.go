package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetClientIPIgnoresForwardedHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	if err := ConfigureTrustedProxies(router, ""); err != nil {
		t.Fatalf("failed to set trusted proxies: %v", err)
	}
	router.GET("/ip", func(c *gin.Context) {
		c.String(http.StatusOK, ClientIPResolver{}.Resolve(c))
	})

	request := httptest.NewRequest(http.MethodGet, "/ip", nil)
	request.RemoteAddr = "203.0.113.7:4321"
	request.Header.Set("X-Forwarded-For", "198.51.100.9")
	request.Header.Set("X-Real-IP", "198.51.100.10")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Body.String() != "203.0.113.7" {
		t.Fatalf("expected 203.0.113.7, got %s", recorder.Body.String())
	}
}

func TestGetClientIPHonorsForwardedHeadersFromTrustedProxy(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	if err := ConfigureTrustedProxies(router, "203.0.113.7/32"); err != nil {
		t.Fatalf("failed to set trusted proxies: %v", err)
	}
	router.GET("/ip", func(c *gin.Context) {
		c.String(http.StatusOK, ClientIPResolver{}.Resolve(c))
	})

	request := httptest.NewRequest(http.MethodGet, "/ip", nil)
	request.RemoteAddr = "203.0.113.7:4321"
	request.Header.Set("X-Forwarded-For", "198.51.100.9")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Body.String() != "198.51.100.9" {
		t.Fatalf("expected 198.51.100.9, got %s", recorder.Body.String())
	}
}

func TestConfigureTrustedProxiesRejectsInvalidValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := ConfigureTrustedProxies(gin.New(), "not-a-cidr")

	if err == nil {
		t.Fatal("expected invalid trusted proxy to fail")
	}
}
