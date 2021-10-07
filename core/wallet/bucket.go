package main

import (
	"encoding/binary"
	"errors"
)

const (
	header       = int(9)
	header64     = int64(header)
	bucketSize   = int(6)
	bucketSize64 = int64(bucketSize)
	//bucketBytes   = int(bucketSize*(size+8) + 8)
	//bucketBytes64 = int64(bucketBytes)
)

var errBucketOverflow = errors.New("bucket overflow")

type Bucket struct {
	n     int64
	data  []byte
	store *BucketStore
}

func (b *Bucket) WriteBulk(data []byte) {
	store := b.store.store
	if int64(len(data)) > b.store.bucketBytes-8 {
		panic(errOverflow)
	}
	offset := header64 + b.n*b.store.bucketBytes
	if _, err := store.Seek(offset, BeginOfFile); err != nil {
		panic(err)
	}
	if n, err := store.Write(data); n != len(data) {
		panic(err)
	}
}

func (b *Bucket) Write(item int64, data []byte) {
	store := b.store.store
	if len(data) != int(b.store.itemBytes) {
		panic("bucket only read entire items")
	}
	if (item+1)*b.store.itemBytes > b.store.bucketBytes-8 {
		panic(errOverflow)
	}
	offset := header64 + b.n*b.store.bucketBytes + item*b.store.itemBytes
	if _, err := store.Seek(offset, BeginOfFile); err != nil {
		panic(err)
	}
	if n, err := store.Write(data); n != len(data) {
		panic(err)
	}
}

func (b *Bucket) Read(item int64) []byte {
	if item > bucketSize64 || item < 0 {
		panic("invalid bucket read")
	}
	return b.data[item*b.store.itemBytes : (item+1)*b.store.itemBytes]
}

func (b *Bucket) GetOverflow() *Bucket {
	overflow := int64(binary.LittleEndian.Uint64(b.data[b.store.size*b.store.itemBytes:]))
	if overflow > b.store.bucketCount {
		panic("bucket overflow")
	}
	if overflow == 0 {
		return nil
	}
	return b.store.ReadBucket(overflow)
}

func (bs *BucketStore) insertBulk(items []byte, bucket int64) {
	b := &Bucket{n: bucket, data: make([]byte, bs.bucketBytes), store: bs}
	bucketCounts := int64(len(items)) / (bs.itemBytes * bucketSize64)
	if bucketCounts == 0 {
		b.WriteBulk(items)
	}
	for n := int64(0); n < bucketCounts; n++ {
		b = b.AddOverflow()
		if n < bucketCounts-1 {
			data := items[n*b.store.bucketBytes : (n+1)*b.store.bucketBytes]
			b.WriteBulk(data)
		} else {
			data := items[(bucketCounts-1)*b.store.bucketBytes:]
			b.WriteBulk(data)
		}
	}
}

func (b *Bucket) ZeroOverflow() {
	b.Write(b.store.itemBytes*b.store.size, make([]byte, 8))
}

func (b *Bucket) AddOverflow() *Bucket {
	data := make([]byte, b.store.bucketBytes+8)
	if _, err := b.store.store.Seek(0, EndOfFile); err != nil {
		panic(err)
	}
	b.store.store.Write(data)
	overflow := b.store.bucketCount
	bucketBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bucketBytes, uint64(overflow))
	b.Write(b.store.itemBytes*b.store.size, bucketBytes)
	b.store.bucketCount++
	return &Bucket{
		n:     overflow,
		data:  data,
		store: b.store,
	}
}

func (b *Bucket) SetOverflow(freeOverFlow int64) {
	if freeOverFlow == -1 {
		data := make([]byte, b.store.bucketBytes)
		if _, err := b.store.store.Seek(0, EndOfFile); err != nil {
			panic(err)
		}
		b.store.store.Write(data)
		freeOverFlow = b.store.bucketCount
		b.store.bucketCount++
	}
	bucketBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bucketBytes, uint64(freeOverFlow))
	b.Write(b.store.itemBytes*b.store.size, bucketBytes)
}

type BucketStore struct {
	store       ByteStore
	bucketCount int64
	itemBytes   int64
	bucketBytes int64
	size        int64
}

func NewBucketStore(itemBytes int64, store ByteStore) *BucketStore {
	if (store.Size()-header64)%(bucketSize64*itemBytes+8) != 0 {
		panic("ByteStore size incompatible with bucket store")
	}
	return &BucketStore{
		store:       store,
		bucketCount: (store.Size() - header64) / (bucketSize64*itemBytes + 8),
		itemBytes:   itemBytes,
		size:        bucketSize64,
		bucketBytes: (bucketSize64*itemBytes + 8),
	}
}

func (b *BucketStore) ReadBucket(n int64) *Bucket {
	if n >= b.bucketCount {
		panic(errBucketOverflow)
	}
	if _, err := b.store.Seek(header64+n*b.bucketBytes, BeginOfFile); err != nil {
		panic(err)
	}
	data := make([]byte, b.bucketBytes)
	if n, err := b.store.Read(data); int64(n) != b.bucketBytes {
		panic(err)
	}
	return &Bucket{
		n:     n,
		data:  data,
		store: b,
	}
}
