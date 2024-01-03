package utils

import "encoding/binary"

func EncodeEndian(value uint32, length int) []byte {
	byteSlice := make([]byte, length)
	binary.LittleEndian.PutUint32(byteSlice, value)
	return byteSlice
}

func DecodeEndian(byteSlice []byte) uint32 {
	buffer := make([]byte, 4)
	copy(buffer, byteSlice)

	return binary.LittleEndian.Uint32(buffer)
}
