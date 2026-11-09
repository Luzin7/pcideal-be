package middlewares

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var ErrInvalidAllowedOrigin = errors.New("invalid CORS_ALLOWED_ORIGINS entry")

func DefaultAllowedOrigins() []string {
	return []string{
		"https://pcideal-web.vercel.app",
	}
}

func ParseAllowedOrigins(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return DefaultAllowedOrigins(), nil
	}

	entries := strings.Split(raw, ",")
	origins := make([]string, 0, len(entries))

	for _, entry := range entries {
		origin, err := normalizeOrigin(entry)
		if err != nil {
			return nil, err
		}

		origins = append(origins, origin)
	}

	return origins, nil
}

func normalizeOrigin(entry string) (string, error) {
	value := strings.TrimSuffix(strings.TrimSpace(entry), "/")
	if value == "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidAllowedOrigin, entry)
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrInvalidAllowedOrigin, entry)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("%w: %q", ErrInvalidAllowedOrigin, entry)
	}

	if parsed.Host == "" || strings.Contains(parsed.Host, "*") {
		return "", fmt.Errorf("%w: %q", ErrInvalidAllowedOrigin, entry)
	}

	if parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", fmt.Errorf("%w: %q", ErrInvalidAllowedOrigin, entry)
	}

	return parsed.Scheme + "://" + parsed.Host, nil
}
