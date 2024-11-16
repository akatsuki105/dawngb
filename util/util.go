package util

// Bit check val's idx bit
func Bit[V uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64](val V, idx int) bool {
	if idx < 0 || idx > 63 {
		return false
	}
	return (val & (1 << idx)) != 0
}

func SetBit[V uint | uint8 | uint16 | uint32 | int8](val V, idx int, b bool) V {
	old := val
	if b {
		val = old | (1 << idx)
	} else {
		val = old & ^(1 << idx)
	}
	return val
}

func Btou8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}
