package middlewares

import (
	"reflect"
	"testing"
)

func TestParseAllowedOriginsUsesSingleProductionOriginWhenUnset(t *testing.T) {
	origins, err := ParseAllowedOrigins("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"https://pcideal-web.vercel.app"}
	if !reflect.DeepEqual(origins, expected) {
		t.Fatalf("expected %v, got %v", expected, origins)
	}
}

func TestParseAllowedOriginsTrimsAndDropsTrailingSlash(t *testing.T) {
	origins, err := ParseAllowedOrigins(" https://pcideal-web.vercel.app/ ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []string{"https://pcideal-web.vercel.app"}
	if !reflect.DeepEqual(origins, expected) {
		t.Fatalf("expected %v, got %v", expected, origins)
	}
}

func TestParseAllowedOriginsRejectsInvalidEntries(t *testing.T) {
	invalid := []string{
		"pcideal-web.vercel.app",
		"https://pcideal-web.vercel.app/app",
		"https://pcideal-web.vercel.app?query=1",
		"*",
		"https://pcideal-web.vercel.app,,https://pcideal-web.vercel.app",
		"ftp://pcideal-web.vercel.app",
	}

	for _, raw := range invalid {
		if _, err := ParseAllowedOrigins(raw); err == nil {
			t.Fatalf("expected %q to be rejected", raw)
		}
	}
}
