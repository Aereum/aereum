package message

func PutByteArray(b []byte, data *[]byte) {
	if len(b) == 0 {
		*data = append(*data, 0)
	}
	if len(b) > 65535 {
		*data = append(*data, append([]byte{255, 255}, b[0:65535]...)...)
	}
	v := len(b)
	*data = append(*data, append([]byte{byte(v), byte(v >> 8)}, b...)...)
}

func PutString(value string, data *[]byte) {
	PutByteArray([]byte(value), data)
}

func PutUint64(v uint64, data *[]byte) {
	b := make([]byte, 8)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	*data = append(*data, b...)
}

func ParseByteArray(data []byte, position int) ([]byte, int) {
	if position+1 >= len(data) {
		return []byte{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	if position+length >= len(data) {
		return []byte{}, position
	}
	return (data[position+1 : position+length+1]), position + length
}

func ParseString(data []byte, position int) (string, int) {
	bytes, newPosition := ParseByteArray(data, position)
	if bytes != nil {
		return string(bytes), newPosition
	} else {
		return "", newPosition
	}
}

func ParseUint64(data []byte, position int) (uint64, int) {
	if position+8 >= len(data) {
		return 0, position + 8
	}
	value := uint64(data[position+0]) |
		uint64(data[position+1])<<8 |
		uint64(data[position+2])<<16 |
		uint64(data[position+3])<<24 |
		uint64(data[position+4])<<32 |
		uint64(data[position+5])<<40 |
		uint64(data[position+6])<<48 |
		uint64(data[position+7])<<56
	return value, position + 8
}
