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
package instructionsold

import (
	"reflect"
	"testing"
)

func TestByteArray(t *testing.T) {
	zero := make([]byte, 0)
	bytes := make([]byte, 0)
	PutByteArray(zero, &bytes)
	inverse, _ := ParseByteArray(bytes, 0)
	if ok := reflect.DeepEqual(zero, inverse); !ok {
		t.Errorf("Wrong ByteArray of zero length")
	}
	large := make([]byte, 256*256+1)
	for n := 0; n < 256*256+1; n++ {
		large[n] = 1
	}
	bytes = make([]byte, 0)
	PutByteArray(large, &bytes)
	inverse, _ = ParseByteArray(bytes, 0)
	if len(inverse) != 1<<16-1 {
		t.Errorf("Wrong ByteArray of large length")
	}
}

func TestString(t *testing.T) {
	bytes := make([]byte, 0)
	PutString("$¢ह€𐍈", &bytes)
	inverse, _ := ParseString(bytes, 0)
	if inverse != "$¢ह€𐍈" {
		t.Errorf("Wrong utf-8 string encoding")
	}
}

func TestUint64(t *testing.T) {
	bytes := make([]byte, 0)
	PutUint64(1<<64-1, &bytes)
	inverse, _ := ParseUint64(bytes, 0)
	if inverse != 1<<64-1 {
		t.Errorf("Wrong uint64 serialization")
	}
}
