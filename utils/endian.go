package utils

import "encoding/binary"

func EncodeEndian(value uint64, length int) []byte {
	byteSlice := make([]byte, length)
	binary.LittleEndian.PutUint64(byteSlice, value)
	return byteSlice
}

func DecodeEndian(byteSlice []byte) uint64 {
	buffer := make([]byte, 8)
	copy(buffer, byteSlice)

	return binary.LittleEndian.Uint64(buffer)
}
