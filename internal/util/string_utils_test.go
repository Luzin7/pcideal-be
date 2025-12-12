package util

import (
	"testing"
)

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "should convert uppercase to lowercase",
			input:    "HELLO WORLD",
			expected: "hello world",
		},
		{
			name:     "should remove special characters",
			input:    "Hello@World#2024!",
			expected: "helloworld2024",
		},
		{
			name:     "should handle mixed case and special chars",
			input:    "Intel Core i7-12700K",
			expected: "intel core i712700k",
		},
		{
			name:     "should preserve spaces between words",
			input:    "AMD Ryzen 7 5800X",
			expected: "amd ryzen 7 5800x",
		},
		{
			name:     "should handle multiple spaces",
			input:    "Hello    World  Test",
			expected: "hello world test",
		},
		{
			name:     "should handle empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "should handle only special characters",
			input:    "@#$%^&*()",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeString(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeString(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected int
	}{
		{
			name:     "identical strings",
			a:        "hello",
			b:        "hello",
			expected: 0,
		},
		{
			name:     "one character difference",
			a:        "hello",
			b:        "hallo",
			expected: 1,
		},
		{
			name:     "completely different strings",
			a:        "abc",
			b:        "xyz",
			expected: 3,
		},
		{
			name:     "empty string vs non-empty",
			a:        "",
			b:        "hello",
			expected: 5,
		},
		{
			name:     "both empty strings",
			a:        "",
			b:        "",
			expected: 0,
		},
		{
			name:     "insertion required",
			a:        "cat",
			b:        "cart",
			expected: 1,
		},
		{
			name:     "deletion required",
			a:        "cart",
			b:        "cat",
			expected: 1,
		},
		{
			name:     "multiple operations",
			a:        "kitten",
			b:        "sitting",
			expected: 3,
		},
		{
			name:     "case sensitive",
			a:        "Hello",
			b:        "hello",
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Levenshtein(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Levenshtein(%q, %q) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestConvertToCents(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{
			name:     "convert 10 reais to cents",
			input:    10,
			expected: 1000,
		},
		{
			name:     "convert 100 reais to cents",
			input:    100,
			expected: 10000,
		},
		{
			name:     "convert 0 reais to cents",
			input:    0,
			expected: 0,
		},
		{
			name:     "convert 1 real to cents",
			input:    1,
			expected: 100,
		},
		{
			name:     "convert large amount",
			input:    5000,
			expected: 500000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCents(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertToCents(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertCentsToReal(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int64
	}{
		{
			name:     "convert 1000 cents to reais",
			input:    1000,
			expected: 10,
		},
		{
			name:     "convert 10000 cents to reais",
			input:    10000,
			expected: 100,
		},
		{
			name:     "convert 0 cents to reais",
			input:    0,
			expected: 0,
		},
		{
			name:     "convert 100 cents to reais",
			input:    100,
			expected: 1,
		},
		{
			name:     "convert large amount",
			input:    500000,
			expected: 5000,
		},
		{
			name:     "convert with rounding (99 cents)",
			input:    99,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCentsToReal(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertCentsToReal(%d) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkNormalizeString(b *testing.B) {
	input := "Intel Core i7-12700K @ 3.6GHz"
	for i := 0; i < b.N; i++ {
		NormalizeString(input)
	}
}

func BenchmarkLevenshtein(b *testing.B) {
	a := "Intel Core i7-12700K"
	b2 := "Intel Core i9-12900K"
	for i := 0; i < b.N; i++ {
		Levenshtein(a, b2)
	}
}
