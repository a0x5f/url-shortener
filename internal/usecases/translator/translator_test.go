package translator

import (
	"log"
	"math"
	"testing"
)

func Test_getShortLinkById(t *testing.T) {
	tests := []struct {
		name     string
		id       uint64
		expected string
	}{
		{
			name:     "Default",
			id:       1,
			expected: ShortLinkDomain + "AAAAAAAAAB",
		},
		{
			name:     "Zero value",
			id:       0,
			expected: ShortLinkDomain + "AAAAAAAAAA",
		},
		{
			name:     "Random number",
			id:       1337,
			expected: ShortLinkDomain + "AAAAAAAAVO",
		},
		{
			name:     "uint64 max",
			id:       math.MaxUint64,
			expected: ShortLinkDomain + "St6VW1zpdiP",
		},
	}

	urlTranslator := UrlTranslator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := urlTranslator.getShortLinkById(tt.id)
			if actual != tt.expected {
				log.Fatalf("expected: %s, actual: %s", tt.expected, actual)
			}
		})
	}
}

func Test_getLinkIdByHash(t *testing.T) {
	tests := []struct {
		name     string
		hash     string
		expected uint64
	}{
		{
			name:     "Default",
			hash:     "AAAAAAAAAB",
			expected: 1,
		},
		{
			name:     "Zero value",
			hash:     "AAAAAAAAAA",
			expected: 0,
		},
		{
			name:     "Random number",
			hash:     "AAAAAAAAVO",
			expected: 1337,
		},
		{
			name:     "uint64 max",
			hash:     "St6VW1zpdiP",
			expected: math.MaxUint64,
		},
	}

	urlTranslator := UrlTranslator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := urlTranslator.getLinkIdByHash(tt.hash)
			if actual != tt.expected {
				log.Fatalf("expected: %d, actual: %d", tt.expected, actual)
			}
		})
	}
}
