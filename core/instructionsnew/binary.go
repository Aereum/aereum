// Copyright 2021 The aereum Authors
// This file is part of the aereum library.
//
// The aereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The aereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package instructionsnew

import "time"

func PutTokenCipher(tc TokenCipher, data *[]byte) {
	PutByteArray(tc.token, data)
	PutByteArray(tc.cipher, data)
}

func PutTokenCiphers(tcs TokenCiphers, data *[]byte) {
	if len(tcs) == 0 {
		*data = append(*data, 0, 0)
		return
	}
	maxLen := len(tcs)
	if len(tcs) > 1<<16-1 {
		maxLen = 1 << 16
	}
	*data = append(*data, byte(maxLen), byte(maxLen>>8))
	for n := 0; n < maxLen; n++ {
		PutTokenCipher(tcs[n], data)
	}
}

func PutByteArray(b []byte, data *[]byte) {
	if len(b) == 0 {
		*data = append(*data, 0, 0)
		return
	}
	if len(b) > 1<<16-1 {
		*data = append(*data, append([]byte{255, 255}, b[0:1<<16-1]...)...)
		return
	}
	v := len(b)
	*data = append(*data, append([]byte{byte(v), byte(v >> 8)}, b...)...)
}

func PutString(value string, data *[]byte) {
	PutByteArray([]byte(value), data)
}

func PutUint16(v uint16, data *[]byte) {
	*data = append(*data, byte(v), byte(v>>8))
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

func PutTime(value time.Time, data *[]byte) {
	bytes, err := value.MarshalBinary()
	if err != nil {
		panic("invalid time")
	}
	PutByteArray(bytes, data)
}

func PutBool(b bool, data *[]byte) {
	if b {
		*data = append(*data, 1)
	} else {
		*data = append(*data, 0)
	}
}

func PutByte(b byte, data *[]byte) {
	*data = append(*data, b)
}

func ParseTokenCipher(data []byte, position int) (TokenCipher, int) {
	tc := TokenCipher{}
	if position+1 >= len(data) {
		return tc, position
	}
	tc.token, position = ParseByteArray(data, position)
	tc.cipher, position = ParseByteArray(data, position)
	return tc, position
}

func ParseTokenCiphers(data []byte, position int) (TokenCiphers, int) {
	if position+1 >= len(data) {
		return TokenCiphers{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	position += 2
	if length == 0 {
		return TokenCiphers{}, position + 2
	}
	if position+length+2 > len(data) {
		return TokenCiphers{}, position + length + 2
	}
	tcs := make(TokenCiphers, length)
	for n := 0; n < length; n++ {
		tcs[n], position = ParseTokenCipher(data, position)
	}
	return tcs, position
}

func ParseByteArray(data []byte, position int) ([]byte, int) {
	if position+1 >= len(data) {
		return []byte{}, position
	}
	length := int(data[position+0]) | int(data[position+1])<<8
	if length == 0 {
		return []byte{}, position + 2
	}
	if position+length+2 > len(data) {
		return []byte{}, position + length + 2
	}
	return (data[position+2 : position+length+2]), position + length + 2
}

func ParseString(data []byte, position int) (string, int) {
	bytes, newPosition := ParseByteArray(data, position)
	if bytes != nil {
		return string(bytes), newPosition
	} else {
		return "", newPosition
	}
}

func ParseUint16(data []byte, position int) (uint16, int) {
	if position+1 >= len(data) {
		return 0, position + 2
	}
	value := uint16(data[position+0]) |
		uint16(data[position+1])<<8
	return value, position + 2
}

func ParseUint64(data []byte, position int) (uint64, int) {
	if position+7 >= len(data) {
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

func ParseTime(data []byte, position int) (time.Time, int) {
	bytes, newposition := ParseByteArray(data, position)
	var t *time.Time
	if err := t.UnmarshalBinary(bytes); err != nil {
		panic("cannot parse time")
	}
	return *t, newposition

}

func ParseBool(data []byte, position int) (bool, int) {
	if position >= len(data) {
		return false, position + 1
	}
	return data[position] != 0, position + 1
}

func ParseByte(data []byte, position int) (byte, int) {
	if position >= len(data) {
		return 0, position + 1
	}
	return data[position], position + 1
}