package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const testRemoteAddr = "203.0.113.7:4321"

func resolveClientIP(t *testing.T, resolver ClientIPResolver, headers map[string]string) string {
	t.Helper()

	gin.SetMode(gin.TestMode)

	router := gin.New()
	if err := ConfigureTrustedProxies(router, ""); err != nil {
		t.Fatalf("failed to set trusted proxies: %v", err)
	}
	router.GET("/ip", func(c *gin.Context) {
		c.String(http.StatusOK, resolver.Resolve(c))
	})

	request := httptest.NewRequest(http.MethodGet, "/ip", nil)
	request.RemoteAddr = testRemoteAddr
	for name, value := range headers {
		request.Header.Set(name, value)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	return recorder.Body.String()
}

func TestClientIPResolverUsesConfiguredHeader(t *testing.T) {
	resolver, err := NewClientIPResolver("CF-Connecting-IP")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := resolveClientIP(t, resolver, map[string]string{"CF-Connecting-IP": "198.51.100.9"})

	if got != "198.51.100.9" {
		t.Fatalf("expected 198.51.100.9, got %s", got)
	}
}

func TestClientIPResolverFallsBackWhenHeaderIsInvalid(t *testing.T) {
	resolver, err := NewClientIPResolver("CF-Connecting-IP")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for _, invalid := range []string{"", "not-an-ip", "198.51.100.9, 203.0.113.10"} {
		got := resolveClientIP(t, resolver, map[string]string{"CF-Connecting-IP": invalid})

		if got != "203.0.113.7" {
			t.Fatalf("expected fallback 203.0.113.7 for %q, got %s", invalid, got)
		}
	}
}

func TestClientIPResolverWithoutHeaderIgnoresHeader(t *testing.T) {
	resolver, err := NewClientIPResolver("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	got := resolveClientIP(t, resolver, map[string]string{"CF-Connecting-IP": "198.51.100.9"})

	if got != "203.0.113.7" {
		t.Fatalf("expected 203.0.113.7, got %s", got)
	}
}

func TestNewClientIPResolverRejectsCommaList(t *testing.T) {
	if _, err := NewClientIPResolver("CF-Connecting-IP, X-Forwarded-For"); err == nil {
		t.Fatal("expected comma list to be rejected")
	}
}
