package index

import (
	"encoding/binary"
	"io"
)

type IO interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
}

func WriteOrPanic(dest IO, bytes []byte, offset int64) {
	if n, _ := dest.WriteAt(bytes, offset); n != len(bytes) {
		panic("IO failure")
	}
}

func ReadOrPanic(src IO, offset int64, size int) []byte {
	bytes := make([]byte, size)
	if n, _ := src.ReadAt(bytes, offset); n != size {
		panic("IO failure")
	}
	return bytes
}

func ByteArrayToUint64Array(bytes []byte) []uint64 {
	count := (len(bytes)) / 8
	data := make([]uint64, 0, count)
	for n := 1; n < count; n++ {
		value := binary.LittleEndian.Uint64(bytes[n*8 : (n+1)*8])
		if n > 0 && value == 0 {
			return data
		}
		data = append(data, value)
	}
	return data
}

func Uint64ArrayToByteArray(data []uint64) []byte {
	bytes := make([]byte, 8*len(data))
	for n := 0; n < len(data); n++ {
		binary.LittleEndian.PutUint64(bytes[n*8:(n+1)*8], data[n])
	}
	return bytes
}
