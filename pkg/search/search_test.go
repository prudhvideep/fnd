package search

import (
	"testing"
)

// TestGetSearchType tests the getSearchType function
func TestGetSearchType(t *testing.T) {
	tests := []struct {
		searchType string
		expected   int
	}{
		{"f", SearchFile},
		{"d", SearchDir},
		{"", SearchNormal},
		{"invalid", SearchNormal},
	}

	for _, tt := range tests {
		t.Run(tt.searchType, func(t *testing.T) {
			result := getSearchType(tt.searchType)
			if result != tt.expected {
				t.Errorf("getSearchType(%s) = %v; want %v", tt.searchType, result, tt.expected)
			}
		})
	}
}
