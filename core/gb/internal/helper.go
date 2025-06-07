package internal

import "golang.org/x/exp/constraints"

func SetByte[V constraints.Unsigned, W constraints.Integer](data V, offset W, value uint8) V {
	return data&^(0xFF<<(int(offset)*8)) | V(value)<<(int(offset)*8)
}
