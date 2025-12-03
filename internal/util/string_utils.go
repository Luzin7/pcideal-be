package util

import (
	"math"
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

func min(a, b, c int) int {
	return int(math.Min(float64(a), math.Min(float64(b), float64(c))))
}

func Levenshtein(a, b string) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	prev := make([]int, len(a)+1)
	curr := make([]int, len(a)+1)

	for i := 0; i <= len(a); i++ {
		prev[i] = i
	}

	for j := 1; j <= len(b); j++ {
		curr[0] = j

		for i := 1; i <= len(a); i++ {
			if a[i-1] == b[j-1] {
				curr[i] = prev[i-1]
			} else {
				curr[i] = min(
					prev[i]+1,
					curr[i-1]+1,
					prev[i-1]+1,
				)
			}
		}

		prev, curr = curr, prev
	}

	return prev[len(a)]
}
