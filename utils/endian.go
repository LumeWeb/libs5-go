package utils

import "encoding/binary"

func EncodeEndian(value uint32, length int) []byte {
	byteSlice := make([]byte, length)
	binary.LittleEndian.PutUint32(byteSlice, value)
	return byteSlice
}

func DecodeEndian(byteSlice []byte) uint32 {
	return binary.LittleEndian.Uint32(byteSlice)
}
