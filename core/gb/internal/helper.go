package internal

import "golang.org/x/exp/constraints"

func Byte[V constraints.Unsigned, W constraints.Integer](data V, offset W) uint8 {
	return uint8(data >> (int(offset) * 8))
}

func SetByte[V constraints.Unsigned, W constraints.Integer](data V, offset W, value uint8) V {
	return data&^(0xFF<<(int(offset)*8)) | V(value)<<(int(offset)*8)
}
