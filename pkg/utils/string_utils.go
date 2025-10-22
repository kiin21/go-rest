package utils

import (
	"strings"
)

// ParseCSVString splits a comma-separated string, trims spaces,
// and removes any empty elements.
func ParseCSVString(input string, sep string) []string {
	// Split by comma
	parts := strings.Split(input, sep)

	// Prepare a clean slice
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
