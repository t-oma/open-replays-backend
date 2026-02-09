// Package parse provides utility functions for parsing various data formats.
package parse

import "fmt"

// FileSizeToBytes parses a human-readable file size string (e.g., "10MB", "1GB").
func FileSizeToBytes(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}

	var multiplier int64 = 1
	switch {
	case len(s) > 2 && (s[len(s)-2:] == "GB" || s[len(s)-2:] == "gb"):
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-2]
	case len(s) > 2 && (s[len(s)-2:] == "MB" || s[len(s)-2:] == "mb"):
		multiplier = 1024 * 1024
		s = s[:len(s)-2]
	case len(s) > 2 && (s[len(s)-2:] == "KB" || s[len(s)-2:] == "kb"):
		multiplier = 1024
		s = s[:len(s)-2]
	}

	var value int64
	if _, err := fmt.Sscanf(s, "%d", &value); err != nil {
		return 0, fmt.Errorf("invalid file size format: %s", s)
	}

	return value * multiplier, nil
}

// FileSizeToHumanReadable parses a file size in bytes to a human-readable string (e.g., "10MB", "1GB").
func FileSizeToHumanReadable(size int64) string {
	if size == 0 {
		return "0B"
	}

	var multiplier int64 = 1
	switch {
	case size > 1024*1024*1024:
		multiplier = 1024 * 1024 * 1024
		size /= multiplier
	case size > 1024*1024:
		multiplier = 1024 * 1024
		size /= multiplier
	case size > 1024:
		multiplier = 1024
		size /= multiplier
	}

	return fmt.Sprintf(
		"%d%s",
		size,
		map[int64]string{1: "B", 1024: "KB", 1024 * 1024: "MB", 1024 * 1024 * 1024: "GB"}[multiplier],
	)
}
