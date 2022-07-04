package main

import (
	"crypto"
	"encoding/binary"
	"io"
)

const (
	BUCKETITEMS = 8
	BUCKETBYTES = (BUCKETITEMS + 1) * 8
)

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

type ReaderWriterAt interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
}

type IndexStore struct {
	hashIndex   map[crypto.Hash]uint64
	io          ReaderWriterAt
	bucketCount uint64
}

func (store *IndexStore) readBucket(bucket uint64) []uint64 {
	bytes := make([]byte, BUCKETBYTES)
	if n, _ := store.io.ReadAt(bytes, int64(bucket)*BUCKETBYTES); n != BUCKETBYTES {
		return nil
	}
	return ByteArrayToUint64Array(bytes)
}

func (store *IndexStore) readBucketSequence(bucket uint64) []uint64 {
	allBuckets := make([]uint64, 0)
	for {
		newBucket := store.readBucket(bucket)
		if newBucket[0] == 0 {
			return append(allBuckets, newBucket[1:int64(newBucket[0])+1]...)
		}
		bucket = newBucket[0]
		allBuckets = append(allBuckets, newBucket[1:]...)
	}
}

func (store *IndexStore) writeOnBucket(bucket uint64, position int, value uint64) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, value)
	offset := int64(bucket)*BUCKETBYTES + int64(position)*8
	store.io.WriteAt(bytes, offset)
}

func (store *IndexStore) appendOnBucket(bucket, value uint64) uint64 {
	data := store.readBucket(bucket)
	if len(data) < BUCKETBYTES+1 {
		store.writeOnBucket(bucket, len(data)+1, value)
		return 0
	}
	return store.appendNewBucket(bucket, value)
}

func (store *IndexStore) appendNewBucket(bucket, value uint64) uint64 {
	store.bucketCount += 1
	bytes := make([]byte, BUCKETBYTES)
	binary.LittleEndian.PutUint64(bytes[0:8], bucket)
	binary.LittleEndian.PutUint64(bytes[8:16], value)
	if n, _ := store.io.WriteAt(bytes, int64(store.bucketCount-1)*BUCKETBYTES); n != BUCKETBYTES {
		panic("could not write new bucket")
	}
	return store.bucketCount
}

func (store *IndexStore) Append(hash crypto.Hash, sequence uint64) {
	store.
}