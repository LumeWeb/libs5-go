package utils

import "encoding/binary"

func EncodeEndian(value uint64, length int) []byte {
	byteSlice := make([]byte, length)
	binary.BigEndian.PutUint64(byteSlice, value)
	return byteSlice
}

func DecodeEndian(byteSlice []byte) uint64 {
	buffer := make([]byte, 8)
	copy(buffer, byteSlice)

	return binary.BigEndian.Uint64(buffer)
}
