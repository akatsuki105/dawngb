package internal

import "golang.org/x/exp/constraints"

// Bit check val's idx bit
func Bit[V constraints.Integer, W constraints.Integer](val V, idx W) bool {
	if idx < 0 || idx > 63 {
		return false
	}
	return (val & (1 << idx)) != 0
}

func SetBit[V constraints.Integer, W constraints.Integer](val V, idx W, b bool) V {
	if b {
		return val | (1 << idx)
	} else {
		return val & ^(1 << idx)
	}
}
