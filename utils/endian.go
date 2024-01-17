package utils

func EncodeEndian(value uint64, length int) []byte {
	res := make([]byte, length)

	for i := 0; i < length; i++ {
		res[i] = byte(value & 0xff)
		value = value >> 8
	}
	return res
}

func DecodeEndian(bytes []byte) uint64 {
	var total uint64

	for i := 0; i < len(bytes); i++ {
		total += uint64(bytes[i]) << (8 * uint(i))
	}

	return total
}
