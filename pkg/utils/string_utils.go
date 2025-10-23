package utils

import (
	"strings"
)

func ParseString(input string, sep string) []string {
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
