package main

import "encoding/binary"

const (
	accBytes   = int(size + 8)
	accBytes64 = int64(accBytes)
)

func CheckHashOnBucket(b *BucketItem, hash Hash, n int) (ok bool, balance int64) {
	data := make([]byte, accBytes)
	b.Read(int64(n)*accBytes64, data)
	check := true
	for n := 0; n < size; n++ {
		if data[n] != hash[n] {
			check = false
			break
		}
	}
	if check {
		return ok, int64(binary.LittleEndian.Uint64(data[bucketSize*accBytes:]))
	}
	return false, 0
}
