// Package datasize helps data size definition.
//
// It is recommended to dot import this package.
package datasize

import "fmt"

const (
	_       = iota
	KB uint = 1 << (10 * iota)
	MB
)

// FormatSize convert 1024 into "1KB"
func FormatSize(s uint) string {
	switch {
	case s < KB:
		return fmt.Sprintf("%dB", s)
	case s < MB:
		return fmt.Sprintf("%dKB", s/KB)
	default:
		return fmt.Sprintf("%dB", s)
	}
}
