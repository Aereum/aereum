package index

import (
	"encoding/binary"
	"os"
)

const (
	bucketItems = 8
	bucketSize  = (bucketItems + 1) * 8
)

type bucket struct {
	storage IO
	count   uint64
}

func (b *bucket) readSingle(bucket uint64) []uint64 {
	bytes := ReadOrPanic(b.storage, int64(bucket)*bucketSize, bucketSize)
	return ByteArrayToUint64Array(bytes)
}

func (b *bucket) readAll(bucket uint64) []uint64 {
	all := make([]uint64, 0)
	for {
		data := b.readSingle(bucket)
		if data[0] == 0 {
			for n := 1; n <= bucketItems; n++ {
				if data[n] == 0 {
					return append(all, data[1:n]...)
				}
			}
			return append(all, data[1:]...)
		}
		bucket = data[0]
	}
}

func (b *bucket) append(bucket, value uint64) uint64 {
	data := b.readSingle(bucket)
	for n := 1; n <= bucketItems; n++ {
		if data[n] == 0 {
			bytes := make([]byte, 8)
			binary.LittleEndian.PutUint64(bytes, value)
			WriteOrPanic(b.storage, bytes, int64(bucket)*bucketSize+int64(n)*8)
			return 0
		}
	}
	bytes := make([]byte, bucketSize)
	binary.LittleEndian.PutUint64(bytes[0:8], bucket)
	binary.LittleEndian.PutUint64(bytes[8:16], value)
	WriteOrPanic(b.storage, bytes, int64(b.count)*bucketSize)
	b.count += 1
	return b.count
}

func (b *bucket) new(value uint64) uint64 {
	bytes := make([]byte, bucketSize)
	binary.LittleEndian.PutUint64(bytes[0:8], 0)
	binary.LittleEndian.PutUint64(bytes[8:16], value)
	WriteOrPanic(b.storage, bytes, int64(b.count)*bucketSize)
	b.count += 1
	return b.count
}

func (b *bucket) Close() {
	b.storage.Close()
}

func OpenBucket(FileName string) *bucket {
	file, err := os.OpenFile(FileName, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return nil
	}
	if stat, err := file.Stat(); err != nil {
		return nil
	} else {
		if stat.Size()%bucketSize != 0 {
			return nil
		}
		return &bucket{
			storage: file,
			count:   uint64(stat.Size() / bucketSize),
		}
	}
}
