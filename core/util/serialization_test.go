package util

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
