package util

import (
	"strings"
	"unicode"
)

func NormalizeString(s string) string {
	s = strings.ToLower(s)

	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsSpace(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(result.String()), " ")
}

func ConvertToCents(price int64) int64 {
	return price * 100
}
