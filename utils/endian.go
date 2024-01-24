package utils

func EncodeEndian(value uint64, length int) []byte {
	res := make([]byte, length)

	for i := length - 1; i >= 0; i-- {
		res[i] = byte(value & 0xff)
		value = value >> 8
	}
	return res
}
func DecodeEndian(bytes []byte) uint64 {
	var total uint64

	for _, b := range bytes {
		total = total*256 + uint64(b)
	}

	return total
}
