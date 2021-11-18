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

// Package message contains data types related to aereum network.
package instruction2

import "time"

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

func PutTime(value time.Time, data *[]byte) {
	bytes, err := value.MarshalBinary()
	if err != nil {
		panic("invalid time")
	}
	PutByteArray(bytes, data)
}

func ParseTime(data []byte, position int) (time.Time, int) {
	bytes, newposition := ParseByteArray(data, position)
	var t *time.Time
	if err := t.UnmarshalBinary(bytes); err != nil {
		panic("cannot parse time")
	}
	return *t, newposition

}
