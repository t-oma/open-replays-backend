// Package stringutil provides utility functions for strings.
package stringutil

import "path/filepath"

// TrimExt removes the file extension from the given filename.
func TrimExt(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}
