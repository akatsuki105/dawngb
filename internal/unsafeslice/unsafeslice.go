package unsafeslice

import "unsafe"

func ByteSliceFromUint16Slice(arr []uint16) []uint8 {
	return unsafe.Slice((*uint8)(unsafe.Pointer(unsafe.SliceData(arr))), len(arr)*2)
}

func ByteSliceFromUint32Slice(arr []uint32) []uint8 {
	return unsafe.Slice((*uint8)(unsafe.Pointer(unsafe.SliceData(arr))), len(arr)*4)
}

func Uint16SliceFromByteSlice(arr []uint8) []uint16 {
	return unsafe.Slice((*uint16)(unsafe.Pointer(unsafe.SliceData(arr))), len(arr)/2)
}

func Uint32SliceFromByteSlice(arr []uint8) []uint32 {
	return unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(arr))), len(arr)/4)
}
